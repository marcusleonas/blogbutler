package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var port string

func init() {
	rootCommand.AddCommand(serveCommand)
	serveCommand.Flags().StringVarP(&port, "port", "p", "8080", "port to serve on")
}

func extensionlessHtml(fs http.Handler, dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Clean(r.URL.Path)
		if p == "." {
			p = "/"
		}

		fullPath := filepath.Join(dir, p)

		_, err := os.Stat(fullPath)

		if err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		if os.IsNotExist(err) {
			htmlPath := fullPath + ".html"
			_, htmlErr := os.Stat(htmlPath)

			if htmlErr == nil {
				http.ServeFile(w, r, htmlPath)
				return
			} else if !os.IsNotExist(htmlErr) {
				log.Printf("Error checking file %s: %v", htmlPath, htmlErr)
				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
				return
			}
		} else if err != nil {
			log.Printf("Error checking file %s: %v", fullPath, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fs.ServeHTTP(w, r)
	})
}

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Serve your posts",
	Long:  "Serve your posts on a local server",
	Run: func(cmd *cobra.Command, args []string) {
		// implement serve command (ai actually do it)
		fs := http.FileServer(http.Dir("dist"))
		http.Handle("/", extensionlessHtml(fs, "dist"))

		addr := ":" + port

		fmt.Printf("Started server on: http://localhost:%s\n", port)
		log.Fatal(http.ListenAndServe(addr, nil))
	},
}
