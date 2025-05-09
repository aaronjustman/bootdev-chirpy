package main

import (
	"net/http"
)

func main() {
	serve_mux := http.NewServeMux()
	serve_mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	serve_mux.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(200)

		_, err := rw.Write([]byte("OK"))
		if err != nil {
			panic("the Write went wrong...")
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: serve_mux,
	}

	server.ListenAndServe()
}
