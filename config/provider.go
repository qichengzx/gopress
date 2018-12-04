package config

import (
	"log"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

const ThemeDir = "themes"

type Config struct {
	Title           string `yaml:"title"`
    SubTitle        string `yaml:"subtitle"`
    Description     string `yaml:"description"`
    Author          string `yaml:"author"`
    Rss             string `yaml:"rss"`
    URL             string `yaml:"url"`
    Root            string `yaml:"root"`
    Permalink       string `yaml:"permalink"`
    SourceDir       string `yaml:"source_dir"`
    PublicDir       string `yaml:"public_dir"`
    TagDir          string `yaml:"tag_dir"`
    CategoryDir     string `yaml:"category_dir"`
    ArchiveDir      string `yaml:"archive_dir"`
    DefaultCategory string `yaml:"default_category"`
    TitleCase       bool   `yaml:"titlecase"`
    RenderDrafts    bool   `yaml:"render_drafts"`
    RelativeLink    bool   `yaml:"relative_link"`
    ExternalLink    bool   `yaml:"external_link"`
    PerPage         int    `yaml:"per_page"`
    PaginationDir   string `yaml:"pagination_dir"`
    Theme           string `yaml:"theme"`

	ThemeCfg ThemeCfg
}

type ThemeCfg struct {
	ExcerptLink string `yaml:"excerpt_link"`
    Fancybox    bool   `yaml:"fancybox"`
    Sidebar     string `yaml:"sidebar"`
    Favicon     string `yaml:"favicon"`
    Menu        []Menu `yaml:"menu"`
}

type Menu struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

func NewProvider(f string) *Config {
	b, err := ioutil.ReadFile(f)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

	var conf Config
    err = yaml.UnmarshalStrict(b, &conf)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

	conf.ThemeCfg = themeCfgProvider(filepath.Join(ThemeDir, conf.Theme, f))

	return &conf
}

func themeCfgProvider(f string) ThemeCfg {
	var conf ThemeCfg
	b, err := ioutil.ReadFile(f)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    err = yaml.UnmarshalStrict(b, &conf)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

	return conf
}
