package main

import (
	"log"
	"net/http"
)

func main() {
	gothmogHandler := &GothmogHandler{}

	mux := http.NewServeMux()
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
