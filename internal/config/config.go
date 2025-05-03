package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	title string
}

func LoadConfig() Config {
	content, err := os.ReadFile("config.toml")
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
