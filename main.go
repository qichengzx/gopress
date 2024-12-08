package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qichengzx/gopress/generator"
)

func main() {
	generatePtr := flag.Bool("g", false, "Generate static files")
	servePtr := flag.Bool("s", false, "Serve static files")
	portPtr := flag.String("p", "8092", "Server port")
	flag.Parse()

	switch {
	case *generatePtr:
		timeStart := time.Now()
		log.Println("Process begining")
		defer func() {
			timeEnd := time.Now()
			used := timeEnd.Sub(timeStart)
			log.Println("Process done")
			log.Println("Used", used)
		}()

		var site = generator.New("./_config.yaml")
		site.Build()
		break

	case *servePtr && *portPtr != "":
		fmt.Println("server running at", *portPtr)
		err := http.ListenAndServe(":"+*portPtr, http.FileServer(http.Dir("./public/")))
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	default:
		flag.Usage()
	}
}
