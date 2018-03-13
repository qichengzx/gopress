package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":8092"

func main() {
	fmt.Println("server run at", port)
	err := http.ListenAndServe(port, http.FileServer(http.Dir("./public/")))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
