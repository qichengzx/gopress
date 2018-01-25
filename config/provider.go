package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Title       string `toml:"title"`
	SubTitle    string `toml:"subtitle"`
	Description string `toml:"description"`

	SourceDir   string `toml:"source_dir"`
	PublicDir   string `toml:"public_dir"`
	TagDir      string `toml:"tag_dir"`
	CategoryDir string `toml:"category_dir"`

	PerPage       int    `toml:"per_page"`
	PaginationDir string `toml:"pagination_dir"`

	ThemeDir string `toml:"theme_dir"`
	Theme    string `toml:"theme"`
}

func NewProvider(f string) *Config {
	var conf Config
	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}

	return &conf
}
