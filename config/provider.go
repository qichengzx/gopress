package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Title       string
	SubTitle    string
	Description string
}

func NewProvider(f string) *Config {
	var conf Config
	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}

	return &conf
}
