package params

import "github.com/iancoleman/strcase"

type EntityParams struct {
	Directory string
	Module    string
	Name      string
	Snake     string
	Camel     string
	Short     string
}

func NewEntityParams(directory, module, name, short string) *EntityParams {
	return &EntityParams{
		Directory: directory,
		Module:    module,
		Name:      name,
		Snake:     strcase.ToSnake(name),
		Camel:     strcase.ToLowerCamel(name),
		Short:     short,
	}
}
