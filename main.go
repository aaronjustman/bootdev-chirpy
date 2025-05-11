package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type api_config struct {
	fileserver_hits atomic.Int32
}

type chirp struct {
	Body string `json:"body"`
}

type chirp_error struct {
	Error string `json:"error"`
}

type chirp_success struct {
	Valid bool `json:"valid"`
}

func (cfg *api_config) increment_fileserver_hits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = cfg.fileserver_hits.Add(1)
		handler.ServeHTTP(w, r)
	})
}

func (cfg *api_config) reset_fileserver_hits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserver_hits.Store(0)

	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	var b []byte
	_, err := w.Write(fmt.Appendf(b, "Hits: %d", cfg.fileserver_hits.Load()))
	if err != nil {
		panic("the Write went wrong...")
	}
}

func (cfg *api_config) write_fileserver_hits(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)

	var b []byte
	html := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserver_hits.Load())
	_, err := w.Write(fmt.Append(b, html))
	if err != nil {
		panic("the Write went wrong...")
	}
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	var chirp chirp
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chirp); err != nil {
		fmt.Println("an error occurred decoding the chirp:", err)
		chirp_error := chirp_error{
			Error: fmt.Sprint("Something went wrong"),
		}
		data, err := json.Marshal(&chirp_error)
		if err != nil {
			fmt.Println("error marshaling the chirp error")
		}
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		w.Write(data)
		return
	}
	defer r.Body.Close()

	if len(chirp.Body) > 140 {
		chirp_error := chirp_error{
			Error: fmt.Sprint("Chirp is too long"),
		}
		data, err := json.Marshal(&chirp_error)
		if err != nil {
			fmt.Println("error marshaling the chirp error")
		}
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	ok_response := chirp_success{
		Valid: true,
	}
	data, err := json.Marshal(&ok_response)
	if err != nil {
		fmt.Println("error marshaling the response string")
	}
	w.Write(data)
}

func main() {
	serve_mux := http.NewServeMux()
	cfg := &api_config{}

	serve_mux.Handle("/app/", cfg.increment_fileserver_hits(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serve_mux.HandleFunc("GET /api/healthz", func(rw http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(200)

		_, err := rw.Write([]byte("OK"))
		if err != nil {
			panic("the Write went wrong...")
		}
	})
	serve_mux.HandleFunc("GET /admin/metrics", cfg.write_fileserver_hits)
	serve_mux.HandleFunc("POST /admin/reset", cfg.reset_fileserver_hits)
	serve_mux.HandleFunc("POST /api/validate_chirp", validate_chirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: serve_mux,
	}

	server.ListenAndServe()
}
