package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"x/config"
	"x/post"
	"x/xlib"
)

var (
	cfFile  string = "./_config.toml"
	appPath string
)

func main() {
	timeStart := time.Now()
	log.Println("Process begining")
	defer func() {
		timeEnd := time.Now()
		used := timeEnd.Sub(timeStart)
		log.Println("Process done")
		log.Println("Used", used)
	}()

	var s xlib.Site
	s.Cfg = config.NewProvider(cfFile)
	themeCfg := filepath.Join(xlib.ThemeDir, s.Cfg.Theme, cfFile)
	s.ThemeCfg = config.ThemeCfgProvider(themeCfg)

	fmt.Println("Welcome to ", s.Cfg.Title)

	appPath, _ := os.Getwd()
	postPath := filepath.Join(appPath, s.Cfg.SourceDir)

	post.Root = s.Cfg.Root
	post.Permalink = s.Cfg.Permalink

	pw, tags, cates := post.GetPosts(postPath)

	tagStr := strings.Join(tags, " ")
	cateStr := strings.Join(cates, " ")

	s.Posts = pw.Posts
	s.CatPosts = pw.CatPosts
	s.TagPosts = pw.TagPosts
	s.Tags = post.WordToMAP(tagStr)
	s.Categories = post.WordToMAP(cateStr)

	s.Build()
}
