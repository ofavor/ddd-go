package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// Create entity command

var entityName string
var shortName string

type entityParams struct {
	ModuleName     string
	EntityName     string
	ShortName      string
	SnakeName      string
	LowerCamelName string
	KebabName      string
}

var eparams entityParams

func init() {
	rootCmd.AddCommand(entityCmd)

	entityCmd.Flags().StringVarP(&moduleName, "module", "m", "", "project module name (required)")
	entityCmd.Flags().StringVarP(&projectDir, "dir", "d", "", "project directory (required)")
	entityCmd.Flags().StringVarP(&entityName, "entity", "e", "", "entity name (required)")
	entityCmd.Flags().StringVarP(&shortName, "short", "s", "", "entity short name (required)")
	entityCmd.MarkFlagRequired("entity")
}

var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Generate entity related files",
	Long:  `Generate entity related files, such as entity, vo, repository, service, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating entity [" + entityName + "]...")
		fmt.Println("Project directory: " + projectDir)
		fmt.Println("Module name: " + moduleName)
		fmt.Println("Entity name: " + entityName)
		fmt.Println("Entity short name: " + shortName)
		eparams = entityParams{moduleName, entityName, shortName, strcase.ToSnake(entityName), strcase.ToLowerCamel(entityName), strcase.ToKebab(entityName)}

		if err := generateCommand(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateDto(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateQuery(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateAppSvc(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateDomainSvc(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateEntity(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateRepo(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateDao(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if err := generatePersist(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateForm(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateView(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := generateRestHandler(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := updateRegistry(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := updateRegistryInfr(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := updateConstant(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := updateConstant(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := updateErrors(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := updateMain(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Success!")
	},
}

func generateCommand() error {
	str := `package command

import "github.com/ofavor/ddd-go/pkg/operator"

type Create{{ .EntityName }}Command struct {
	Operator *operator.UserOperator
	// TODO
}

type Update{{ .EntityName }}Command struct {
	Operator *operator.UserOperator
	Id int64
	// TODO
}

type Delete{{ .EntityName }}Command struct {
	Operator *operator.UserOperator
	Id int64
	// TODO
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/application/command/%s_cmd.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateDto() error {
	str := `package dto

import "{{ .ModuleName }}/internal/domain/model"

type {{ .EntityName }}Dto struct {
	Id int64
	// TODO
}

func To{{ .EntityName }}Dto(m model.{{ .EntityName }}) *{{ .EntityName }}Dto {
	return &{{ .EntityName }}Dto{
		Id: m.GetId(),
		// TODO
	}
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/application/dto/%s_dto.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateQuery() error {
	str := `package query

import "github.com/ofavor/ddd-go/pkg/operator"

type Get{{ .EntityName }}ByIdQuery struct {
	Operator *operator.UserOperator
	Id int64
	// TODO
}

type Get{{ .EntityName }}ListQuery struct {
	Operator *operator.UserOperator
	Offset int64
	Limit  int64
	// TODO
}

	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/application/query/%s_qry.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateAppSvc() error {
	str := `package application

import (
	{{ .ModuleName }} "{{ .ModuleName }}/internal"
	"{{ .ModuleName }}/internal/application/command"
	"{{ .ModuleName }}/internal/application/dto"
	"{{ .ModuleName }}/internal/application/query"
)

type {{ .EntityName }}ApplicationService interface {
	Create{{ .EntityName }}(cmd *command.Create{{ .EntityName }}Command) (int64, error)
	Update{{ .EntityName }}(cmd *command.Update{{ .EntityName }}Command) error
	Delete{{ .EntityName }}(cmd *command.Delete{{ .EntityName }}Command) error

	Get{{ .EntityName }}ById(query *query.Get{{ .EntityName }}ByIdQuery) (*dto.{{ .EntityName }}Dto, error)
	Get{{ .EntityName }}List(query *query.Get{{ .EntityName }}ListQuery) (int64, []*dto.{{ .EntityName }}Dto, error)
}

type {{ .LowerCamelName }}ApplicationService struct {
	registry {{ .ModuleName }}.Registry
}

func New{{ .EntityName }}ApplicationService() {{ .EntityName }}ApplicationService {
	reg := {{ .ModuleName }}.GetRegistry()
	return &{{ .LowerCamelName }}ApplicationService{
		registry: reg,
	}
}

// Create{{ .EntityName }} implements {{ .EntityName }}ApplicationService.
func (s *{{ .LowerCamelName }}ApplicationService) Create{{ .EntityName }}(cmd *command.Create{{ .EntityName }}Command) (int64, error) {
	panic("unimplemented")
}

// Update{{ .EntityName }} implements {{ .EntityName }}ApplicationService.
func (s *{{ .LowerCamelName }}ApplicationService) Update{{ .EntityName }}(cmd *command.Update{{ .EntityName }}Command) error {
	panic("unimplemented")
}

// Delete{{ .EntityName }} implements {{ .EntityName }}ApplicationService.
func (s *{{ .LowerCamelName }}ApplicationService) Delete{{ .EntityName }}(cmd *command.Delete{{ .EntityName }}Command) error {
	panic("unimplemented")
}

// Get{{ .EntityName }}ById implements {{ .EntityName }}ApplicationService.
func (s *{{ .LowerCamelName }}ApplicationService) Get{{ .EntityName }}ById(query *query.Get{{ .EntityName }}ByIdQuery) (*dto.{{ .EntityName }}Dto, error) {
	panic("unimplemented")
}

// Get{{ .EntityName }}List implements {{ .EntityName }}ApplicationService.
func (s *{{ .LowerCamelName }}ApplicationService) Get{{ .EntityName }}List(query *query.Get{{ .EntityName }}ListQuery) (int64, []*dto.{{ .EntityName }}Dto, error) {
	panic("unimplemented")
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/application/%s_app_svc.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateDomainSvc() error {
	str := `package service

type {{ .EntityName }}DomainService interface {
	// TODO
}

type {{ .LowerCamelName }}DomainService struct {
	// TODO
}

func New{{ .EntityName }}DomainService() {{ .EntityName }}DomainService {
	return &{{ .LowerCamelName }}DomainService{
		// TODO
	}
}
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/service/%s_svc.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}
func generateEntity() error {
	str := `package model

import (
	"github.com/ofavor/ddd-go/pkg/idgen"
	"{{ .ModuleName }}/internal/domain/constant"
)

type {{ .EntityName }}Id = int64

// {{ .EntityName }} entity interface
type {{ .EntityName }} interface {
	GetId() {{ .EntityName }}Id
	GetCreatedAt() int64
	GetUpdatedAt() int64
}

type {{ .LowerCamelName }} struct {
	id        {{ .EntityName }}Id
	createdAt int64
	updatedAt int64
}

func New{{ .EntityName }}(
	idgen idgen.IdGenerator,
) ({{ .EntityName }}, error) {
	id, err := idgen.NextId(constant.Name{{ .EntityName }})
	if err != nil {
		return nil, err
	}
	return &{{ .LowerCamelName }}{
		id: id,
	}, nil
}

func Load{{ .EntityName }}(
	id int64,
) {{ .EntityName }} {
	return &{{ .LowerCamelName }}{
		id: id,
	}
}

// GetId implements {{ .EntityName }}.
func (e *{{ .LowerCamelName }}) GetId() {{ .EntityName }}Id {
	return e.id
}

// GetCreatedAt implements {{ .EntityName }}.
func (e *{{ .LowerCamelName }}) GetCreatedAt() int64 {
	return e.createdAt
}

// GetUpdatedAt implements {{ .EntityName }}.
func (e *{{ .LowerCamelName }}) GetUpdatedAt() int64 {
	return e.updatedAt
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/model/%s.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateRepo() error {
	str := `package repository

import (
	"github.com/ofavor/ddd-go/pkg/repo"
	"{{ .ModuleName }}/internal/domain/model"
)

type {{ .EntityName }}Filter struct {
	Id []int64
}

func (f *{{ .EntityName }}Filter) Conditions() map[string]interface{} {
	conds := make(map[string]interface{})
	if len(f.Id) > 0 {
		conds["id IN (?)"] = f.Id
	}
	// Add custom conditions here
	// TODO
	return conds
}

type {{ .EntityName }}Repository interface {
	repo.Repository[model.{{ .EntityName }}]
	// Add custom methods here
	// TODO
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/domain/repository/%s_repo.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateDao() error {
	str := `package dao

	import (
		"{{ .ModuleName }}/internal/domain/constant"
		"{{ .ModuleName }}/internal/domain/model"
	)
	
	type {{ .EntityName }}Dao struct {
		Id int64 ` + "`gorm:\"autoIncrement:false;primaryKey\"`" + `
	}
	
	func (o {{ .EntityName }}Dao) TableName() string {
		return constant.Name{{ .EntityName }}
	}
	
	func (o {{ .EntityName }}Dao) SequenceName() string {
		return constant.Name{{ .EntityName }}
	}
	
	func (o {{ .EntityName }}Dao) InitialId() int64 {
		return 10000
	}
	
	func To{{ .EntityName }}Dao(m model.{{ .EntityName }}) *{{ .EntityName }}Dao {
		return &{{ .EntityName }}Dao{
			Id: m.GetId(),
		}
	}
	
	func To{{ .EntityName }}Entity(d *{{ .EntityName }}Dao) model.{{ .EntityName }} {
		return model.Load{{ .EntityName }}(
			d.Id,
		)
	}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/infrustructure/repository/dao/%s_dao.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generatePersist() error {
	str := `package persist

import (
	repogorm "github.com/ofavor/ddd-go/pkg/repo/gorm"
	"{{ .ModuleName }}/internal/domain/model"
	"{{ .ModuleName }}/internal/domain/repository"
	"{{ .ModuleName }}/internal/infrustructure/repository/dao"

	"gorm.io/gorm"
)

type {{ .LowerCamelName }}Repository struct {
	*repogorm.GormRepo[model.{{ .EntityName }}, dao.{{ .EntityName }}Dao]
}

func New{{ .EntityName }}Repository(db *gorm.DB) repository.{{ .EntityName }}Repository {
	return &{{ .LowerCamelName }}Repository{
		GormRepo: repogorm.NewRepo(
			db,
			func(d *dao.{{ .EntityName }}Dao) model.{{ .EntityName }} {
				return dao.To{{ .EntityName }}Entity(d)
			},
			func(e model.{{ .EntityName }}) *dao.{{ .EntityName }}Dao {
				return dao.To{{ .EntityName }}Dao(e)
			},
		),
	}
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/infrustructure/repository/persist/%s_repo.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateForm() error {
	str := `package form

	type Create{{ .EntityName }}Form struct {
	}
	
	type Update{{ .EntityName }}Form struct {
	}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/interfaces/rest/form/%s_form.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateView() error {
	str := `package view

import "{{ .ModuleName }}/internal/application/dto"

type {{ .EntityName }}Vm struct {
	Id int64 ` + "`json:\"id\"`" + `
}

func New{{ .EntityName }}Vm(o *dto.{{ .EntityName }}Dto) *{{ .EntityName }}Vm {
	return &{{ .EntityName }}Vm{
		Id: o.Id,
	}
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/interfaces/rest/view/%s_vm.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func generateRestHandler() error {
	str := `package rest

import (
	"github.com/ofavor/ddd-go/pkg/operator"
	"github.com/ofavor/ddd-go/pkg/rest"
	"{{ .ModuleName }}/internal/application"
	"{{ .ModuleName }}/internal/application/command"
	"{{ .ModuleName }}/internal/application/query"
	"{{ .ModuleName }}/internal/interfaces/rest/form"
	"{{ .ModuleName }}/internal/interfaces/rest/view"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

type {{ .EntityName }}Handler struct {
	{{ .ShortName }}AppSvc application.{{ .EntityName }}ApplicationService
}

func New{{ .EntityName }}Handler() rest.Handler {
	return &{{ .EntityName }}Handler{
		{{ .ShortName }}AppSvc: application.New{{ .EntityName }}ApplicationService(),
	}
}

func (h *{{ .EntityName }}Handler) Name() string {
	return "{{ .LowerCamelName }}"
}

func (h *{{ .EntityName }}Handler) BuildRoutes(ws *restful.WebService) {
	type {{ .LowerCamelName }}Result struct {
		rest.Result
		Data *view.{{ .EntityName }}Vm ` + "`json:\"data\" description:\"{{ .EntityName }} info\"`" + `
	}
	type {{ .LowerCamelName }}IdResult struct {
		rest.Result
		Data int64 ` + "`json:\"data\" description:\"{{ .EntityName }} id\"`" + `
	}
	type {{ .LowerCamelName }}List struct {
		Total int64          ` + "`json:\"total\"`" + `
		List  []*view.{{ .EntityName }}Vm ` + "`json:\"list\" description:\"{{ .EntityName }} list\"`" + `
	}
	type {{ .LowerCamelName }}ListResult struct {
		rest.Result
		Data {{ .LowerCamelName }}List ` + "`json:\"data\" description:\"{{ .EntityName }} list data\"`" + `
	}
	ws.Route(rest.WithAuthHeader(
		ws.GET("/{{ .KebabName }}/{id}").
			To(h.get{{ .EntityName }}ById).
			Doc("Get {{ .LowerCamelName }} by id").
			Metadata(restfulspec.KeyOpenAPITags, []string{"{{ .LowerCamelName }}"}).
			Param(ws.PathParameter("id", "{{ .EntityName }} id")).
			Returns(200, "OK", {{ .LowerCamelName }}Result{}),
	))
	ws.Route(rest.WithAuthHeader(
		ws.GET("/{{ .KebabName }}").
			To(h.get{{ .EntityName }}List).
			Doc("Get {{ .LowerCamelName }} list").
			Metadata(restfulspec.KeyOpenAPITags, []string{"{{ .LowerCamelName }}"}).
			Param(ws.QueryParameter("page", "page number").Required(false).DefaultValue("1")).
			Param(ws.QueryParameter("page_size", "page size").Required(false).DefaultValue("20")).
			Returns(200, "OK", {{ .LowerCamelName }}ListResult{}),
	))
	ws.Route(rest.WithAuthHeader(
		ws.POST("/{{ .KebabName }}").
			To(h.create{{ .EntityName }}).
			Doc("Create {{ .LowerCamelName }}").
			Metadata(restfulspec.KeyOpenAPITags, []string{"{{ .LowerCamelName }}"}).
			Reads(form.Create{{ .EntityName }}Form{}).
			Returns(200, "OK", {{ .LowerCamelName }}IdResult{}),
	))
	ws.Route(rest.WithAuthHeader(
		ws.PUT("/{{ .KebabName }}/{id}").
			To(h.update{{ .EntityName }}).
			Doc("Update {{ .LowerCamelName }}").
			Metadata(restfulspec.KeyOpenAPITags, []string{"{{ .LowerCamelName }}"}).
			Param(ws.PathParameter("id", "{{ .LowerCamelName }} id")).
			Reads(form.Update{{ .EntityName }}Form{}).
			Returns(200, "OK", rest.Result{}),
	))
	ws.Route(rest.WithAuthHeader(
		ws.DELETE("/{{ .KebabName }}/{id}").
			To(h.delete{{ .EntityName }}).
			Doc("Delete {{ .LowerCamelName }}").
			Metadata(restfulspec.KeyOpenAPITags, []string{"{{ .LowerCamelName }}"}).
			Param(ws.PathParameter("id", "{{ .LowerCamelName }} id")).
			Returns(200, "OK", rest.Result{}),
	))
}

func (h *{{ .EntityName }}Handler) get{{ .EntityName }}ById(request *restful.Request, response *restful.Response) {
	id := rest.PathParamAsInt(request, "id", 0)
	{{ .ShortName }}, err := h.{{ .ShortName }}AppSvc.Get{{ .EntityName }}ById(&query.Get{{ .EntityName }}ByIdQuery{Id: id})
	if err != nil {
		rest.Error(response, err)
		return
	}
	rest.Success(response, view.New{{ .EntityName }}Vm({{ .ShortName }}))
}

func (h *{{ .EntityName }}Handler) get{{ .EntityName }}List(request *restful.Request, response *restful.Response) {
	page := rest.QueryParamAsInt(request, "page", 1)
	pageSize := rest.QueryParamAsInt(request, "page_size", 20)
	total, {{ .ShortName }}List, err := h.
		{{ .ShortName }}AppSvc.
		Get{{ .EntityName }}List(&query.Get{{ .EntityName }}ListQuery{
			Operator: operator.NewUserOperatorFromRest(request),
			Offset:   (page - 1) * pageSize,
			Limit:    pageSize,
		})
	if err != nil {
		rest.Error(response, err)
		return
	}
	list := make([]*view.{{ .EntityName }}Vm, 0, len({{ .ShortName }}List))
	for _, u := range {{ .ShortName }}List {
		list = append(list, view.New{{ .EntityName }}Vm(u))
	}
	rest.Success(response, rest.ListData{Total: total, List: list})
}

func (h *{{ .EntityName }}Handler) create{{ .EntityName }}(request *restful.Request, response *restful.Response) {
	form := &form.Create{{ .EntityName }}Form{}
	if err := request.ReadEntity(form); err != nil {
		rest.Error(response, err)
		return
	}
	id, err := h.{{ .ShortName }}AppSvc.Create{{ .EntityName }}(&command.Create{{ .EntityName }}Command{
		Operator: operator.NewUserOperatorFromRest(request),
		// TODO
	})
	if err != nil {
		rest.Error(response, err)
		return
	}
	rest.Success(response, id)
}

func (h *{{ .EntityName }}Handler) update{{ .EntityName }}(request *restful.Request, response *restful.Response) {
	id := rest.PathParamAsInt(request, "id", 0)
	form := &form.Update{{ .EntityName }}Form{}
	if err := request.ReadEntity(form); err != nil {
		rest.Error(response, err)
		return
	}
	err := h.{{ .ShortName }}AppSvc.Update{{ .EntityName }}(&command.Update{{ .EntityName }}Command{
		Operator: operator.NewUserOperatorFromRest(request),
		Id:       id,
		// TODO
	})
	if err != nil {
		rest.Error(response, err)
		return
	}
	rest.Success(response, nil)
}

func (h *{{ .EntityName }}Handler) delete{{ .EntityName }}(request *restful.Request, response *restful.Response) {
	id := rest.PathParamAsInt(request, "id", 0)
	err := h.{{ .ShortName }}AppSvc.Delete{{ .EntityName }}(&command.Delete{{ .EntityName }}Command{
		Operator: operator.NewUserOperatorFromRest(request),
		Id:       id,
	})
	if err != nil {
		rest.Error(response, err)
		return
	}
	rest.Success(response, nil)
}
	
	`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return writeTemplateToFile(fmt.Sprintf("%s/internal/interfaces/rest/%s_handler.go", projectDir, eparams.SnakeName), tpl, eparams)
	}
}

func updateGoSource(path string, strs []string) error {
	fmt.Printf("Updating file: %s ...", path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(raw)

	templateToString := func(src string) (string, error) {
		tpl, err := template.New("").Parse(src)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, eparams); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	for i, str := range strs {
		needle, err := templateToString(str)
		if err != nil {
			return err
		}
		slot := fmt.Sprintf("// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE #%d\n", i+1)
		if !strings.Contains(content, needle) {
			content = strings.ReplaceAll(
				content,
				slot,
				needle+"\n"+slot,
			)
		}
	}
	if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	fmt.Println("        Done")
	return nil
}

func updateRegistry() error {
	return updateGoSource(fmt.Sprintf("%s/internal/registry.go", projectDir), []string{
		`"{{ .ModuleName }}/internal/domain/repository"`,
		`Get{{ .EntityName }}Repository() repository.{{ .EntityName }}Repository`,
	})
}

func updateRegistryInfr() error {
	return updateGoSource(fmt.Sprintf("%s/internal/infrustructure/registry/registry.go", projectDir), []string{
		`"{{ .ModuleName }}/internal/domain/repository"`,
		`"{{ .ModuleName }}/internal/infrustructure/repository/persist"`,
		`"{{ .ModuleName }}/internal/infrustructure/repository/dao"`,
		`{{ .LowerCamelName }}Repo repository.{{ .EntityName }}Repository`,
		`&dao.{{ .EntityName }}Dao{},`,
		`func (r *registry) Get{{ .EntityName }}Repository() repository.{{ .EntityName }}Repository {
	if r.{{ .LowerCamelName }}Repo == nil {
		r.{{ .LowerCamelName }}Repo = persist.New{{ .EntityName }}Repository(r.db.GetConn().(*gorm.DB))
	}
	return r.{{ .LowerCamelName }}Repo
}
		`,
	})
}

func updateConstant() error {
	return updateGoSource(fmt.Sprintf("%s/internal/domain/constant/constant.go", projectDir), []string{
		`Name{{ .EntityName }} = "{{ .ModuleName }}_{{ .SnakeName }}"`,
	})
}

func updateErrors() error {
	return updateGoSource(fmt.Sprintf("%s/internal/domain/errors/errors.go", projectDir), []string{
		`Err{{ .EntityName }}NotFound = errors.NewError(1, "{{ .EntityName }} not found")`,
	})
}

func updateMain() error {
	return updateGoSource(fmt.Sprintf("%s/cmd/main.go", projectDir), []string{
		`irest "{{ .ModuleName }}/internal/interfaces/rest"`,
		`irest.New{{ .EntityName }}Handler(),`,
	})
}
