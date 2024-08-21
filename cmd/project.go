package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// Create project command

var moduleName string
var projectDir string

type projectParams struct {
	ModuleName string
	ProjectDir string
}

var pparams projectParams

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.Flags().StringVarP(&moduleName, "module", "m", "", "project module name (required)")
	projectCmd.Flags().StringVarP(&projectDir, "dir", "d", "", "project directory (required)")
	projectCmd.MarkFlagRequired("name")
	projectCmd.MarkFlagRequired("dir")
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create project",
	Long:  `Create project with all needed DDD directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new project [" + moduleName + "]...")
		fmt.Println("Project directory: " + projectDir)
		pparams = projectParams{moduleName, projectDir}

		if err := prepareDirectories(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareGitFiles(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareGoModFiles(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareEnvFile(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareMakeFile(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareReadmeFile(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareVscodeLaunchFile(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := prepareGoSources(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Success!")
	},
}

func prepareDirectories() error {
	dirs := []string{
		"",
		"cmd",
		"internal",
		"internal/application",
		"internal/application/dto",
		"internal/application/command",
		"internal/application/query",
		"internal/domain",
		"internal/domain/constant",
		"internal/domain/errors",
		"internal/domain/event",
		"internal/domain/model",
		"internal/domain/model/vo",
		"internal/domain/service",
		"internal/domain/repository",
		"internal/infrustructure",
		"internal/infrustructure/registry",
		"internal/infrustructure/repository",
		"internal/infrustructure/repository/dao",
		"internal/infrustructure/repository/persist",
		"internal/interfaces",
		"internal/interfaces/rest",
		"internal/interfaces/rest/form",
		"internal/interfaces/rest/view",
		"internal/interfaces/event",
		".vscode",
	}
	for _, d := range dirs {
		fmt.Printf("Creating directory: %s ...", projectDir+"/"+d)
		if err := os.Mkdir(projectDir+"/"+d, os.ModePerm); err != nil {
			return err
		}
		fmt.Println("        Done")
	}
	return nil
}

func prepareGitFiles() error {
	// return os.WriteFile(projectDir+"/.gitignore", []byte(`
	str := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Build directory
build/

# Dependency directories (remove the comment below to include it)
# vendor/

.DS_Store

# swagger-ui
swagger-ui/

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/.gitignore", projectDir), tpl, pparams)
	}
}

func prepareGoModFiles() error {
	str := `module {{ .ModuleName }}

go 1.21.4

require (
	ddd-go v0.0.0
)

require (
)

replace ddd-go => ../ddd-go
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/go.mod", projectDir), tpl, pparams)
	}
}

func prepareEnvFile() error {
	str := `# Log
LOG_LEVEL=debug

# Server
PORT=8081
SWAGGER_DIR=../swagger-ui/dist
PPROF=true
JWT_KEY=i88U2kmkhwq29dkDD2ybb

# Database
DB_TYPE=mysql
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=root
MYSQL_DBNAME={{ .ModuleName }}
# POSTGRES_HOST=
# POSTGRES_PORT=
# POSTGRES_USER=
# POSTGRES_PASSWORD=
# POSTGRES_DBNAME={{ .ModuleName }}

# Cache
CACHE_TYPE=redis
CACHE_PREFIX={{ .ModuleName }}:
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
	
# Event bus
EVENT_BUS_TYPE=redis
EVENT_CONSUME_GROUP={{ .ModuleName }}-consumer-group	
EVENT_BUFFER_SIZE=1000

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/.env", projectDir), tpl, pparams)
	}
}

func prepareMakeFile() error {
	str := `BUILD_ENV := CGO_ENABLED=0
TARGET_EXEC := {{ .ModuleName }}
BUILD=` + "`date +%FT%T%z`" + `
LDFLAGS=-ldflags "-w -s -X main.AppName=${TARGET_EXEC} -X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: all clean setup build-linux build-osx build-windows

all: clean setup build-linux build-osx build-windows

clean:
	rm -rf build

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

build-linux: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${TARGET_EXEC} cmd/main.go

build-osx: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${TARGET_EXEC} cmd/main.go

build-windows: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${TARGET_EXEC}.exe cmd/main.go
	

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/Makefile", projectDir), tpl, pparams)
	}
}

func prepareReadmeFile() error {
	str := `# {{ .ModuleName }}

## Get Started

Run go command

` + "```bash\ngo run cmd/main.go -config .env\n```" + `

Swagger UI

` + "```\nhttp://127.0.0.1:8081/apidocs\n```" + `
	
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/README.md", projectDir), tpl, pparams)
	}
}

func prepareVscodeLaunchFile() error {
	str := `{
		// Use IntelliSense to learn about possible attributes.
		// Hover to view descriptions of existing attributes.
		// For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
		"version": "0.2.0",
		"configurations": [
			{
				"name": "Launch Main",
				"type": "go",
				"request": "launch",
				"mode": "auto",
				"cwd": "${workspaceFolder}",
				"program": "${workspaceFolder}/cmd"
			}
		]
	}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/.vscode/launch.json", projectDir), tpl, pparams)
	}
}

func generateGoRegistry() error {
	str := `package {{ .ModuleName }}

import (
	"github.com/ofavor/ddd-go/pkg/cache"
	"github.com/ofavor/ddd-go/pkg/event"
	"github.com/ofavor/ddd-go/pkg/idgen"
	"github.com/ofavor/ddd-go/pkg/mutex"
	"github.com/ofavor/ddd-go/pkg/tx"

	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #1
)

type Registry interface {
	GetTransMgr() tx.TransMgr
	GetEventBus() event.EventBus
	GetCache() cache.Cache
	GetMutex() mutex.Mutex
	GetIdGenerator() idgen.IdGenerator

	// repositories
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #2
}

var _registry Registry

func GetRegistry() Registry {
	if _registry == nil {
		panic("registry is not initialized")
	}
	return _registry
}

func OnRegistryInitialized(reg Registry) {
	if _registry != nil {
		panic("registry is already initialized")
	}
	_registry = reg
}	
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/registry.go", projectDir), tpl, pparams)
	}
}

func generateGoMain() error {
	str := `package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ofavor/ddd-go/pkg/log"
	"github.com/ofavor/ddd-go/pkg/rest"
	"github.com/ofavor/ddd-go/pkg/util"
	"{{ .ModuleName }}/internal/infrustructure/registry"

	"github.com/joho/godotenv"
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #1
)

var (
	AppName string = "{{ .ModuleName }}"
	Version string = "0.0.0"
	Build   string = "0"
)

var (
	version *bool
	config  *string
)

func main() {
	version = flag.Bool("version", false, "print version information and exit")
	config = flag.String("config", "", "config file path")
	flag.Parse()

	if *version {
		fmt.Println(AppName + " " + Version + "(" + Build + ")")
		os.Exit(0)
	}

	fmt.Println(AppName + " " + Version + "(" + Build + ") is starting ...")

	if len(*config) > 0 {
		if err := godotenv.Load(*config); err != nil {
			log.Fatalf("Error loading %s file: %s", *config, err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
	}
	log.SetLevel(os.Getenv("LOG_LEVEL"))
	util.JwtKey = os.Getenv("JWT_KEY")

	registry.NewRegistry()

	gw := rest.NewGateway(
		"{{ .ModuleName }}",
		fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")),
		[]rest.Handler{
			// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #2
		},
		os.Getenv("SWAGGER_DIR"),
		os.Getenv("PPROF") == "true",
	)
	go gw.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	s := <-c
	log.Info("Exit with signal ", s)
}
	

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/cmd/main.go", projectDir), tpl, pparams)
	}
}

func generateGoRegistryInfra() error {
	str := `package registry

import (
	"github.com/ofavor/ddd-go/pkg/cache"
	cacheredis "github.com/ofavor/ddd-go/pkg/cache/redis"
	"github.com/ofavor/ddd-go/pkg/db"
	dbgorm "github.com/ofavor/ddd-go/pkg/db/gorm"
	"github.com/ofavor/ddd-go/pkg/event"
	evtredis "github.com/ofavor/ddd-go/pkg/event/redis"
	"github.com/ofavor/ddd-go/pkg/idgen"
	idgengorm "github.com/ofavor/ddd-go/pkg/idgen/gorm"
	"github.com/ofavor/ddd-go/pkg/log"
	"github.com/ofavor/ddd-go/pkg/mutex"
	mutexredis "github.com/ofavor/ddd-go/pkg/mutex/redis"
	"github.com/ofavor/ddd-go/pkg/tx"
	txgorm "github.com/ofavor/ddd-go/pkg/tx/gorm"
	"fmt"
	{{ .ModuleName }} "{{ .ModuleName }}/internal"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #1
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #2
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #3
)

type registry struct {
	db       db.Database
	transMgr tx.TransMgr
	idGen    idgen.IdGenerator
	cache    cache.Cache
	mutex    mutex.Mutex
	evtBus   event.EventBus

	// repos
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #4
}

func NewRegistry() {{ .ModuleName }}.Registry {
	reg := &registry{}
	reg.initDatabase()
	reg.initCache()
	reg.initMutex()
	reg.initIdGenerator()
	reg.initTransManager()
	reg.initEventBus()

	reg.initModels()

	{{ .ModuleName }}.OnRegistryInitialized(reg)

	// initialize data
	go reg.initData()

	return reg
}

func (r *registry) initDatabase() {
	log.Info("Initialize database ...")
	t := os.Getenv("DB_TYPE")
	switch t {
	case "mysql":
		r.db = dbgorm.NewDatabase(
			os.Getenv("DB_TYPE"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DBNAME"),
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
		)
	case "postgres":
		r.db = dbgorm.NewDatabase(
			os.Getenv("DB_TYPE"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DBNAME"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		)
	default:
		panic("Unsupported database type: " + t)
	}
}

func (r *registry) initCache() {
	log.Info("Initialize cache ...")
	t := os.Getenv("CACHE_TYPE")
	switch t {
	case "redis":
		r.cache = cacheredis.NewCache(
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
			os.Getenv("REDIS_PASSWORD"),
			os.Getenv("REDIS_DB"),
			os.Getenv("CACHE_PREFIX"),
		)
	default:
		panic("Unsupported cache type: " + t)
	}
}

func (r *registry) initMutex() {
	log.Info("Initialize mutex ...")
	r.mutex = mutexredis.NewMutex(r.cache.GetConn().(*redis.Client))
}

func (r *registry) initIdGenerator() {
	log.Info("Initialize id generator ...")
	r.idGen = idgengorm.NewIdGenerator(r.db.GetConn().(*gorm.DB))
}

func (r *registry) initTransManager() {
	log.Info("Initialize transaction ...")
	r.transMgr = txgorm.NewTransMgr(r.db.GetConn().(*gorm.DB))
}

func (r *registry) initEventBus() {
	log.Info("Initialize event bus ...")
	t := os.Getenv("EVENT_BUS_TYPE")
	switch t {
	case "redis":
		bufferSize, _ := strconv.ParseInt(os.Getenv("EVENT_BUFFER_SIZE"), 10, 64)
		r.evtBus = evtredis.NewEventBusWithConn(
			r.cache.GetConn().(*redis.Client),
			bufferSize,
			os.Getenv("EVENT_CONSUME_GROUP"),
		)
	case "kafka":
		// TODO
	default:
		panic("Unsupported event bus type: " + t)
	}
}

func (r *registry) initModels() {
	models := []interface{}{
		// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #5
	}

	// initialize database tables ...
	r.db.RegisterModels(models)

	// initialize id generator
	for _, m := range models {
		if e, ok := m.(idgen.IdSupport); ok {
			if err := r.idGen.Initialize(e); err != nil {
				panic(fmt.Errorf("failed to initialize Id generator types:%v", err))
			}
		}
	}
}

func (r *registry) initData() {
}

func (r *registry) GetEventBus() event.EventBus {
	return r.evtBus
}

func (r *registry) GetTransMgr() tx.TransMgr {
	return r.transMgr
}

func (r *registry) GetCache() cache.Cache {
	return r.cache
}

func (r *registry) GetIdGenerator() idgen.IdGenerator {
	return r.idGen
}

func (r *registry) GetMutex() mutex.Mutex {
	return r.mutex
}	
	
// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #6

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/infrustructure/registry/registry.go", projectDir), tpl, pparams)
	}
}

func generateGoConstant() error {
	str := `package constant

var (
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #1
)

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/constant/constant.go", projectDir), tpl, pparams)
	}
}

func generateGoErrors() error {
	str := `package errors

import "github.com/ofavor/ddd-go/pkg/errors"

var (
	ErrPermissionDenied = errors.NewError(1, "permission denied")
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #1
)

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/errors/errors.go", projectDir), tpl, pparams)
	}
}

func generateGoEvent() error {
	str := `package event

var (
	// TODO
)

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/event/event.go", projectDir), tpl, pparams)
	}
}

func prepareGoSources() error {
	if err := generateGoMain(); err != nil {
		return err
	}
	if err := generateGoRegistry(); err != nil {
		return err
	}
	if err := generateGoRegistryInfra(); err != nil {
		return err
	}
	if err := generateGoConstant(); err != nil {
		return err
	}
	if err := generateGoErrors(); err != nil {
		return err
	}
	if err := generateGoEvent(); err != nil {
		return err
	}
	return nil
}
