package cmd

import (
	"fmt"
	"os"

	"github.com/anthonynsimon/new/pkg"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "new [template path] [destination path]",
	Short:   "render custom templates",
	Example: "new templates/team-project .",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a base template path")
			os.Exit(1)
			return
		}
		template := args[0]
		destinationDir := "."
		if len(args) == 2 {
			destinationDir = args[1]
		}

		tmpl := lib.NewTemplate(template, destinationDir)

		if err := tmpl.Resolve(); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		if err := tmpl.Render(); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
