package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/marcusleonas/blogbutler/internal/config"
	"github.com/marcusleonas/blogbutler/internal/utils"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

func init() {
	rootCommand.AddCommand(buildCommand)
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Build your posts to html",
	Long:  "Build all posts to html",
	Run: func(cmd *cobra.Command, args []string) {
		confExists := config.ConfigExists()
		if !confExists {
			fmt.Println("No blogbutler config file found. Run `blogbutler init` to create one.")
			os.Exit(1)
		}

		config := config.LoadConfig()

		_, err := os.Stat(config.Site.PostFolder)
		if err != nil {
			fmt.Println("No posts folder found. Run `blogbutler init` to create one.")
			os.Exit(1)
		}

		_, err = os.Stat("dist")
		if err == nil {
			nErr := os.RemoveAll("dist")
			if nErr != nil {
				fmt.Println(nErr)
				os.Exit(1)
			}
		}

		err = os.Mkdir("dist", 0755)
		if err != nil {
			fmt.Println("Failed to create dist folder.")
			os.Exit(1)
		}

		err = os.Mkdir("dist/posts", 0755)
		if err != nil {
			fmt.Println("Failed to create dist folder.")
			os.Exit(1)
		}

		md := goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				&frontmatter.Extender{},
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)

		var postsWithMeta []utils.Post

		posts, _ := os.ReadDir("posts")
		for _, post := range posts {
			if post.IsDir() || !strings.HasSuffix(post.Name(), ".md") || post.Name() == "index.md" {
				continue
			}

			fmt.Printf("Building post '%s'...\n", post.Name())

			postFilePath := path.Join("posts", post.Name())
			file, err := os.ReadFile(postFilePath)
			if err != nil {
				log.Printf("Error reading post file %s: %v. Skipping.", postFilePath, err)
				continue
			}

			ctx := parser.NewContext()

			var buf bytes.Buffer
			if err := md.Convert(file, &buf, parser.WithContext(ctx)); err != nil {
				log.Printf("Error converting markdown: %v. Skipping.", postFilePath)
				continue
			}
			htmlOutput := buf.Bytes()

			f := frontmatter.Get(ctx)
			var meta struct {
				Title       string `yaml:"title"`
				Description string `yaml:"description"`
				Date        string `yaml:"date"`
				Author      string `yaml:"author"`
			}

			err = f.Decode(&meta)
			if err != nil {
				fmt.Printf("Error decoding frontmatter in '%s'. Required frontmatter probably does not exist.", postFilePath)
				continue
			}

			layoutTemplate := "templates/layout.html"
			postTemplate := "templates/post.html"

			outputFilename := strings.Trim(post.Name(), ".md")

			tmpl, err := template.ParseFiles(layoutTemplate, postTemplate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			data := struct {
				SiteTitle   string
				PostTitle   string
				PostContent template.HTML
				Copyright   string
			}{
				SiteTitle:   config.Site.Title,
				PostTitle:   meta.Title,
				PostContent: template.HTML(htmlOutput),
				Copyright:   config.Site.Copyright,
			}

			outFile, err := os.Create(path.Join("dist", "posts", outputFilename+".html"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = tmpl.ExecuteTemplate(outFile, "layout", data)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer outFile.Close()

			postsWithMeta = append(postsWithMeta, utils.Post{
				PostTitle:   meta.Title,
				PostPath:    "/posts/" + outputFilename,
				Author:      meta.Author,
				Date:        meta.Date,
				Description: meta.Description,
			})

			log.Printf("Successfully built post '%s'.\n", post.Name())
		}

		err = os.Mkdir("dist/public", 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// render index page
		indexOutputFilename := "index.html"

		tmpl, err := template.ParseFiles("templates/layout.html", "templates/home.html")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data := struct {
			PostTitle string
			SiteTitle string
			Posts     []utils.Post
			Copyright string
		}{
			PostTitle: "Home",
			SiteTitle: config.Site.Title,
			Posts:     postsWithMeta,
			Copyright: config.Site.Copyright,
		}

		outFile, err := os.Create(path.Join("dist", indexOutputFilename))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tmpl.ExecuteTemplate(outFile, "layout", data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer outFile.Close()

		log.Printf("Successfully built index page '%s'.\n", indexOutputFilename)

		// copy public assets
		_, err = os.Stat("public")
		if err == nil {
			err = utils.CopyDirectory("public", "dist")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		rssFeedXML, err := utils.GenerateRSS(postsWithMeta, "", config.Site.Title, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating RSS feed: %v\n", err)
			return
		}
		rssOutputFilename := "rss.xml"
		rssOutputFilepath := path.Join("dist", rssOutputFilename)
		err = os.WriteFile(rssOutputFilepath, []byte(rssFeedXML), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing RSS feed to file: %v\n", err)
			return
		}

		log.Println("Successfully built all posts.")
	},
}
