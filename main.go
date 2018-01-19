package main

import (
	"fmt"
	"log"
	"time"
	"x/config"
	"x/post"
)

var cfFile string = "./_config.toml"

type Site struct {
	Posts []*post.Post
	cfg   *config.Config
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

	posts := post.GetPosts("./source/_posts")

	for _, p := range posts {
		fmt.Println(p)
	}
}
