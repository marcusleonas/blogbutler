package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/marcusleonas/blogbutler/internal/config"
	"github.com/marcusleonas/blogbutler/internal/template"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(initCommand)
}

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "initialise a new blogbutler site",
	Long:  "initialise a new blogbutler site",
	Run: func(cmd *cobra.Command, args []string) {
		// check if more than 1 arg
		if len(args) > 1 {
			fmt.Println("Too many arguments")
			os.Exit(1)
		}

		folderName := args[0]
		postFolder := "posts"
		templatesFolder := "templates"

		err := os.MkdirAll(path.Join(folderName, postFolder), 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.Mkdir(path.Join(folderName, templatesFolder), 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = config.CreateConfig(folderName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = template.CreateTemplates(folderName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Successfully initialised in `%s`", folderName)
	},
}
