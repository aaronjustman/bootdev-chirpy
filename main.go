package main

import (
	"net/http"
)

func main() {
	serve_mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: serve_mux,
	}

	server.ListenAndServe()
}
