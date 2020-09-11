package main

import (
	"log"
	"net/http"
)

func main() {
	port := "3000"
	directory := "./test/assets"
	fileServer := http.FileServer(http.Dir(directory))
	http.Handle("/", fileServer)
	log.Printf("Serving %s on HTTP port: %s\n", directory, port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
