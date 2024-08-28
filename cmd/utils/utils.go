package utils

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func WriteTemplateToFile(path string, tpl *template.Template, params interface{}) error {
	fmt.Printf("Creating file: %s ...", path)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("\n%s already exists", path)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := tpl.Execute(f, params); err != nil {
		return err
	}
	fmt.Println("        Done")
	return nil
}

func UpdateGoSource(path string, replaces map[string]string) error {
	fmt.Printf("Updating file: %s ...", path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(raw)

	for key, str := range replaces {
		content = strings.ReplaceAll(
			content,
			key,
			str+"\n"+key,
		)
	}
	if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	fmt.Println("        Done")
	return nil
}
