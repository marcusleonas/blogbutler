package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type SiteConfig struct {
	Title      string `toml:"title"`
	PostFolder string `toml:"post_folder"`
	Copyright  string `toml:"copyright"`
}

type Config struct {
	Site SiteConfig `toml:"site"`
}

const configFileName = "config.toml"
const defaultConfigContent = `# docs: https://github.com/marcusleonas/blogbutler/wiki
[site]
title = "My Default Blog Title"
post_folder = "posts"
copyright = "© 2025 Your Name Here"
`

func CreateConfig(folder string) error {
	fullConfigPath := path.Join(folder, configFileName)

	_, err := os.Stat(fullConfigPath)

	if err == nil {
		fmt.Printf("Config file '%s' already exists.\n", configFileName)
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(
			"error checking config file '%s': %w",
			configFileName,
			err,
		)
	}

	fmt.Printf(
		"Config file '%s' not found. Creating default config...\n",
		configFileName,
	)
	err = os.WriteFile(fullConfigPath, []byte(defaultConfigContent), 0644)
	if err != nil {
		return fmt.Errorf(
			"failed to create default config file '%s': %w",
			configFileName,
			err,
		)
	}

	fmt.Printf("Successfully created default config file '%s'.\n", configFileName)
	return nil
}

func LoadConfig() Config {
	content, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var conf Config
	_, err = toml.Decode(string(content), &conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return conf
}
