package config

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
)

type Config struct {
	Title       string `toml:"title"`
	SubTitle    string `toml:"subtitle"`
	Description string `toml:"description"`
	Author      string `toml:"author"`
	Rss         string `toml:"rss"`

	Url       string `toml:"url"`
	Root      string `toml:"root"`
	Permalink string `toml:"permalink"`

	SourceDir   string `toml:"source_dir"`
	PublicDir   string `toml:"public_dir"`
	TagDir      string `toml:"tag_dir"`
	CategoryDir string `toml:"category_dir"`
	ArchiveDir  string `toml:"archive_dir"`

	DefaultCategory string `toml:"default_category"`

	TitleCase    bool `toml:"titlecase"`
	RelativeLink bool `toml:"relative_link"`

	PerPage       int    `toml:"per_page"`
	PaginationDir string `toml:"pagination_dir"`

	ThemeDir string `toml:"theme_dir"`
	Theme    string `toml:"theme"`
	ThemeCfg ThemeCfg
}

type ThemeCfg struct {
	Menu        []Menu `toml:"menu"`
	ExcerptLink string `toml:"excerpt_link"`
	Fancybox    bool   `toml:"fancybox"`
	Sidebar     string `toml:"sidebar"`
	Favicon     string `toml:"favicon"`
}

type Menu struct {
	Title string `toml:"title"`
	Url   string `toml:"url"`
}

func NewProvider(f string) *Config {
	var conf Config
	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}

	conf.ThemeCfg = themeCfgProvider(filepath.Join(conf.ThemeDir, conf.Theme, f))

	return &conf
}

func themeCfgProvider(f string) ThemeCfg {
	var conf ThemeCfg
	if _, err := toml.DecodeFile(f, &conf); err != nil {
		panic(err)
	}

	return conf
}
