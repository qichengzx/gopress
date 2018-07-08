package main

import (
	"fmt"
	"github.com/qichengzx/gopress/xlib"
	"log"
	"time"
)

var cfFile = "./_config.toml"

func main() {
	timeStart := time.Now()
	log.Println("Process begining")
	defer func() {
		timeEnd := time.Now()
		used := timeEnd.Sub(timeStart)
		log.Println("Process done")
		log.Println("Used", used)
	}()

	var site = xlib.New(cfFile)
	fmt.Println("Welcome to ", site.Cfg.Title)
	site.Build()
}
