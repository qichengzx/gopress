package main

import (
	"flag"
	"fmt"
	"github.com/qichengzx/gopress/generator"
	"log"
	"net/http"
	"os"
	"time"
)

var cfFile = "./_config.yaml"

func main() {
	generatePtr := flag.Bool("g", true, "Generate static files")
	servePtr := flag.Bool("s", false, "Serve static files")
	portPtr := flag.String("p", "8092", "Server port")
	flag.Parse()

	if *generatePtr {
		timeStart := time.Now()
		log.Println("Process begining")
		defer func() {
			timeEnd := time.Now()
			used := timeEnd.Sub(timeStart)
			log.Println("Process done")
			log.Println("Used", used)
		}()

		var site = generator.New(cfFile)
		site.Build()
		os.Exit(0)
	}

	if *servePtr && *portPtr != "" {
		fmt.Println("server run at", *portPtr)
		err := http.ListenAndServe(*portPtr, http.FileServer(http.Dir("./public/")))
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
