package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ddd-go",
	Short: "A generator golang DDD development",
	Long: `ddd-go is a library for golang DDD development. 
It provides some useful tools for DDD development. 
It could also help users generating entity, etc. `,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
