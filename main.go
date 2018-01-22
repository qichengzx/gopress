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
)

var (
	cfFile  string = "./_config.toml"
	appPath string
)

type Site struct {
	Posts      []post.Post
	Categories map[string]int
	Tags       map[string]int

	cfg *config.Config
}

func main() {
	timeStart := time.Now()
	log.Println("Process begining")
	defer func() {
		timeEnd := time.Now()
		used := timeEnd.Sub(timeStart)
		log.Println("Process done")
		log.Println("Used", used)
	}()

	var s Site
	s.cfg = config.NewProvider(cfFile)
	fmt.Println("Welcome to ", s.cfg.Title)

	appPath, _ := os.Getwd()
	postPath := filepath.Join(appPath, s.cfg.SourceDir)

	posts, tags, cates := post.GetPosts(postPath)

	tagStr := strings.Join(tags, " ")
	cateStr := strings.Join(cates, " ")

	s.Posts = posts
	s.Tags = post.WordToMAP(tagStr)
	s.Categories = post.WordToMAP(cateStr)

	fmt.Println(s)
}
