package cmd

import (
	"fmt"
	"os"
	"path"

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
		postFolder := path.Join(folderName, "posts")
		err := os.MkdirAll(postFolder, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		f, err := os.Create(path.Join(folderName, "config.toml"))
		defer f.Close()

		f.WriteString(fmt.Sprintf(`# docs: https://github.com/marcusleonas/blogbutler/wiki
[site]
title="%s"
post_folder="posts"`, folderName))

		fmt.Printf("Successfully initialised in `%s`", folderName)
	},
}
