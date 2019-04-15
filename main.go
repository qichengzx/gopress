package main

import (
	"fmt"
	"github.com/qichengzx/gopress/generator"
	"log"
	"time"
)

var cfFile = "./_config.yaml"

func main() {
	timeStart := time.Now()
	log.Println("Process begining")
	defer func() {
		timeEnd := time.Now()
		used := timeEnd.Sub(timeStart)
		log.Println("Process done")
		log.Println("Used", used)
	}()

	var site = generator.New(cfFile)
	fmt.Println("Welcome to ", site.Cfg.Title)
	site.Build()
}
