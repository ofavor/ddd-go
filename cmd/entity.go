package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ofavor/ddd-go/cmd/kratos"
	"github.com/ofavor/ddd-go/cmd/params"
	"github.com/spf13/cobra"
)

var directory string = "."
var layout string = "kratos"
var module string
var entityName string
var entityShort string

func init() {
	rootCmd.AddCommand(entityCmd)
	entityCmd.Flags().StringVarP(&directory, "directory", "d", directory, "project directory, default is '.'")
	entityCmd.Flags().StringVarP(&layout, "layout", "l", layout, "project layout, default is 'kratos'")
	entityCmd.Flags().StringVarP(&module, "module", "m", module, "project module name")
	entityCmd.MarkFlagRequired("module")
	entityCmd.Flags().StringVarP(&entityName, "name", "n", entityName, "entity name")
	entityCmd.MarkFlagRequired("name")
	entityCmd.Flags().StringVarP(&entityShort, "short", "s", entityShort, "entity short name")
	entityCmd.MarkFlagRequired("short")
}

var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Generate entity related files",
	Long:  `Generate entity related files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Trying to generate entity...")
		params := params.NewEntityParams(directory, module, entityName, entityShort)
		fmt.Println("Entity name:", params.Name)
		fmt.Println("Entity snake name:", params.Snake)
		fmt.Println("Entity camel name:", params.Camel)
		fmt.Println("Entity short name:", params.Short)

		// prompt user to type 'Y' to continue
		fmt.Print("Continue? [Y/n] ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToUpper(confirm) != "Y" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		// check if directory exists
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			fmt.Println("Directory does not exist")
			os.Exit(1)
		}
		switch layout {
		case "kratos":
			kratos.GenerateFiles(params)
		default:
			fmt.Println("Invalid layout")
			os.Exit(1)
		}
	},
}
