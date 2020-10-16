package main

import (
    "net/http"
)

type GothmogHandler struct {
}

func (b *GothmogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    rw.Header().Set("Content-Type", "text/plain")
    rw.Write([]byte ("here we go"))
}
