package kratos

import (
	"fmt"
	"text/template"

	"github.com/ofavor/ddd-go/cmd/params"
	"github.com/ofavor/ddd-go/cmd/utils"
)

func GenerateFiles(params *params.EntityParams) {
	if err := generateEntity(params); err != nil {
		panic(err)
	}
	if err := generateDao(params); err != nil {
		panic(err)
	}
	if err := generateRepository(params); err != nil {
		panic(err)
	}
	if err := generateRepo(params); err != nil {
		panic(err)
	}
	if err := updateInfrastructure(params); err != nil {
		panic(err)
	}
}

func generateEntity(params *params.EntityParams) error {
	str := `package entity

import (
	"{{ .Module }}/internal/infrastructure/repo/dao"

	"github.com/ofavor/ddd-go/pkg/entity"
)

type {{ .Name }}Id = uint

type {{ .Name }} interface {
	entity.Entity[dao.{{ .Name }}Dao]
	GetId() {{ .Name }}Id
	GetName() string
}

type {{ .Camel }} struct {
	dao *dao.{{ .Name }}Dao
}

func New{{ .Name }}(name string) ({{ .Name }}, error) {
	m := &{{ .Camel }}{
		dao: &dao.{{ .Name }}Dao{
			Name: name,
		},
	}
	return m, nil
}

func Load{{ .Name }}(data *dao.{{ .Name }}Dao) {{ .Name }} {
	return &{{ .Camel }}{
		dao: data,
	}
}

func (e *{{ .Camel }}) IsNew() bool {
	return e.dao.ID == 0
}

func (e *{{ .Camel }}) DAO() *dao.{{ .Name }}Dao {
	return e.dao
}

func (e *{{ .Camel }}) GetId() {{ .Name }}Id {
	return e.dao.ID
}

func (e *{{ .Camel }}) GetName() string {
	return e.dao.Name
}
`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return utils.WriteTemplateToFile(fmt.Sprintf("%s/internal/domain/entity/%s.go", params.Directory, params.Snake), tpl, params)
	}
}

func generateDao(params *params.EntityParams) error {
	str := `package dao

import "gorm.io/gorm"

type {{ .Name }}Dao struct {
	gorm.Model
	Name string ` + "`" + `gorm:"type:varchar(50);not null;default:'';comment:Name"` + "`" + `
}

func (d *{{ .Name }}Dao) TableName() string {
	return "prefix_{{ .Snake }}"
}

// func (d *{{ .Name }}Dao) BeforeSave(tx *gorm.DB) error {
// 	return nil
// }

// func (d *{{ .Name }}Dao) AfterFind(tx *gorm.DB) error {
// 	return nil
// }
`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return utils.WriteTemplateToFile(fmt.Sprintf("%s/internal/infrastructure/repo/dao/%s_dao.go", params.Directory, params.Snake), tpl, params)
	}
}

func generateRepository(params *params.EntityParams) error {
	str := `package repository

import (
	"{{ .Module }}/internal/domain/entity"
	"{{ .Module }}/internal/infrastructure/repo/dao"

	"github.com/ofavor/ddd-go/pkg/repo"
)

type {{ .Name }}Filter struct {
	Id       []int64
	Name     string
	NameLike string
}

func (f {{ .Name }}Filter) Conditions() map[string]interface{} {
	conds := make(map[string]interface{})
	if len(f.Id) > 0 {
		conds["id"] = f.Id
	}
	if f.Name != "" {
		conds["name"] = f.Name
	}
	if f.NameLike != "" {
		conds["name like"] = f.NameLike
	}
	return conds
}

type {{ .Name }}Repo interface {
	repo.Repository[entity.{{ .Name }}, dao.{{ .Name }}Dao]
}
`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return utils.WriteTemplateToFile(fmt.Sprintf("%s/internal/domain/repository/%s_repo.go", params.Directory, params.Snake), tpl, params)
	}
}

func generateRepo(params *params.EntityParams) error {
	str := `package repo

import (
	"{{ .Module }}/internal/domain/entity"
	"{{ .Module }}/internal/domain/repository"
	"{{ .Module }}/internal/infrastructure/repo/dao"

	"github.com/ofavor/ddd-go/pkg/db"
	repogorm "github.com/ofavor/ddd-go/pkg/repo/gorm"
	"gorm.io/gorm"
)

type {{ .Camel }}Repo struct {
	repogorm.GormRepo[entity.{{ .Name }}, dao.{{ .Name }}Dao]
}

func New{{ .Name }}Repo(db db.Database) repository.{{ .Name }}Repo {
	return &{{ .Camel }}Repo{
		GormRepo: *repogorm.NewRepo(
			db.GetConn().(*gorm.DB),
			func(d *dao.{{ .Name }}Dao) entity.{{ .Name }} {
				return entity.Load{{ .Name }}(d)
			},
		),
	}
}
`
	if tpl, err := template.New("").Parse(str); err != nil {
		return err
	} else {
		return utils.WriteTemplateToFile(fmt.Sprintf("%s/internal/infrastructure/repo/%s_repo.go", params.Directory, params.Snake), tpl, params)
	}
}

func updateInfrastructure(params *params.EntityParams) error {
	return utils.UpdateGoSource(fmt.Sprintf("%s/internal/infrastructure/infra.go", params.Directory), map[string]string{
		"	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE new repo": fmt.Sprintf("	repo.New%sRepo,", params.Name),
		"		// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE new dao": fmt.Sprintf("		&dao.%sDao{},", params.Name),
	})
}
