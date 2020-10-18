package main

import (
	"log"
	"net/http"
)

func main() {
	gothmogHandler := &GothmogHandler{}

	mux := http.NewServeMux()
	// We don't have any fixed endpoints; each update request will send data in a specific format
	// that gothmogHandler will parse into usable data.
	mux.Handle("/", gothmogHandler)

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
