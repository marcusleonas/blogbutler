package template

import (
	_ "embed"
	"os"
	"path"
)

//go:embed home.html
var HomeTemplate string

//go:embed post.html
var PostsTemplate string

//go:embed layout.html
var LayoutTemplate string

func CreateTemplates(folderName string) error {
	homeTemplate, err := os.Create(path.Join(folderName, "templates", "home.html"))
	if err != nil {
		return err
	}
	defer homeTemplate.Close()
	homeTemplate.WriteString(HomeTemplate)

	layoutTemplate, err := os.Create(path.Join(folderName, "templates", "layout.html"))
	if err != nil {
		return err
	}
	defer layoutTemplate.Close()
	layoutTemplate.WriteString(LayoutTemplate)

	PostTemplate, err := os.Create(path.Join(folderName, "templates", "post.html"))
	if err != nil {
		return err
	}
	defer PostTemplate.Close()
	PostTemplate.WriteString(PostsTemplate)

	IndexTemplate, err := os.Create(path.Join(folderName, "templates", "index.html"))
	if err != nil {
		return err
	}
	defer IndexTemplate.Close()
	IndexTemplate.WriteString(HomeTemplate)
	return nil
}
