package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

func writeTemplateToFile(path string, tpl *template.Template, params interface{}) error {
	fmt.Printf("Creating file: %s ...", path)
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

var rootCmd = &cobra.Command{
	Use:   "ddd-go",
	Short: "A generator golang DDD development",
	Long: `ddd-go is a library for golang DDD development. 
It provides some useful tools for DDD development. 
It could also help users generating projects, models, etc. `,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
