package cmd

import (
	"fmt"
	"os"

	"github.com/marcusleonas/blogbutler/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(buildCommand)
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Build your posts to html",
	Long:  "Build all posts to html",
	Run: func(cmd *cobra.Command, args []string) {
		exists := config.ConfigExists()
		if !exists {
			fmt.Println("No blogbutler config file found. Run `blogbutler init` to create one.")
			os.Exit(1)
		}

		fmt.Println("build command ran")
	},
}
