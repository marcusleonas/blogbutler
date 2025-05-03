package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "blogbutler",
	Short: "Little static blog generator.",
	Long: `
Little static blog generator. Converts markdown files to HTML.
Also provides a server if required.`,
}

func Execute() {
	err := rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
