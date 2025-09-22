package main

import (
	"github.com/kornetvba/YandShortUrl/internal/handler"
	"net/http"
)

func run() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.TextPlainPage)
	mux.HandleFunc("/{id}", handler.GetTextPlainPage)

	return http.ListenAndServe(":8080", mux)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
