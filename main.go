package main

import (
	"net/http"
)

func main() {
	serve_mux := http.NewServeMux()
	serve_mux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr:    ":8080",
		Handler: serve_mux,
	}

	server.ListenAndServe()
}
