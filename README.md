# DDD-GO

Golang DDD 实践框架，可以快速创建DDD项目以及实体相关代码，使开发更简单，更快捷。

## Usage

```bash
$ ddd-go -h
ddd-go is a library for golang DDD development. 
It provides some useful tools for DDD development. 
It could also help users generating projects, models, etc.

Usage:
  ddd-go [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  entity      Generate entity related files
  help        Help about any command
  project     Create project

Flags:
  -h, --help   help for ddd-go

Use "ddd-go [command] --help" for more information about a command.
```

创建项目

```bash
ddd-go project -m myproj -d ../myproj
```

创建实体

```bash
ddd-go entity -m myproj -d ../myproj -e User -s usr
```

创建的项目目录结构

```bash
.
├── cmd
│   └── main.go
├── go.mod
├── .env
└── internal
    ├── application
    │   ├── command
    │   │   └── user_cmd.go
    │   ├── dto
    │   │   └── user_dto.go
    │   ├── query
    │   │   └── user_qry.go
    │   └── user_app_svc.go
    ├── domain
    │   ├── constant
    │   │   └── constant.go
    │   ├── errors
    │   │   └── errors.go
    │   ├── event
    │   │   └── event.go
    │   ├── model
    │   │   ├── user.go
    │   │   └── vo
    │   ├── repository
    │   │   └── user_repo.go
    │   └── service
    │       └── user_svc.go
    ├── infrustructure
    │   ├── registry
    │   │   └── registry.go
    │   └── repository
    │       ├── dao
    │       │   └── user_dao.go
    │       └── persist
    │           └── user_repo.go
    ├── interfaces
    │   ├── event
    │   └── rest
    │       ├── form
    │       │   └── user_form.go
    │       ├── user_handler.go
    │       └── view
    │           └── user_vm.go
    └── registry.go

```

修改.env中的相关参数
```bash
# Log
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
MYSQL_DBNAME=myproj
# POSTGRE_HOST=
# POSTGRE_PORT=
# POSTGRE_USER=
# POSTGRE_PASSWORD=
# POSTGRE_DBNAME=myproj

# Cache
CACHE_TYPE=redis
CACHE_PREFIX=myproj
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
	
# Event bus
EVENT_BUS_TYPE=redis
EVENT_CONSUME_GROUP=myproj-consumer-group	
EVENT_BUFFER_SIZE=1000
```

启动服务

```bash
go run cmd/main.go
```