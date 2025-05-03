package template

import (
	"os"
	"path"
)

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
	return nil
}
