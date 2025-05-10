package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type api_config struct {
	fileserver_hits atomic.Int32
}

func (cfg *api_config) increment_fileserver_hits(handler http.Handler) http.Handler {
	cfg.fileserver_hits.Add(1)
	return handler
}

func (cfg *api_config) reset_fileserver_hits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserver_hits.Store(0)

	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	//w.WriteHeader(fmt.Sprintf("Hits: %d", cfg.fileserver_hits))

	var b []byte
	_, err := w.Write(fmt.Appendf(b, "Hits: %d", cfg.fileserver_hits))
	if err != nil {
		panic("the Write went wrong...")
	}
}

func (cfg *api_config) write_fileserver_hits(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	//w.WriteHeader(fmt.Sprintf("Hits: %d", cfg.fileserver_hits))

	var b []byte
	_, err := w.Write(fmt.Appendf(b, "Hits: %d", cfg.fileserver_hits))
	if err != nil {
		panic("the Write went wrong...")
	}
}

func main() {
	serve_mux := http.NewServeMux()
	cfg := &api_config{}
	serve_mux.Handle("/app/", cfg.increment_fileserver_hits(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serve_mux.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(200)

		_, err := rw.Write([]byte("OK"))
		if err != nil {
			panic("the Write went wrong...")
		}
	})
	serve_mux.HandleFunc("/metrics", cfg.write_fileserver_hits)
	serve_mux.HandleFunc("/reset", cfg.reset_fileserver_hits)

	server := http.Server{
		Addr:    ":8080",
		Handler: serve_mux,
	}

	server.ListenAndServe()
}
