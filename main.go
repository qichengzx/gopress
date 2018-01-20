package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"x/config"
	"x/post"
)

var (
	cfFile   string = "./_config.toml"
	postPath string = "./source/_posts"
)

type Site struct {
	Posts      []*post.Post
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

	posts, tags, cates := post.GetPosts("./source/_posts")

	tagStr := strings.Join(tags, " ")
	cateStr := strings.Join(cates, " ")

	s.Posts = posts
	s.Tags = post.WordToMAP(tagStr)
	s.Categories = post.WordToMAP(cateStr)

	fmt.Println(s)
}
