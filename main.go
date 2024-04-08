package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	address = ":8080"
	timeout = 30 * time.Second
)

func main() {
	// wire up dependencies
	w := make([]int64, 0)
	m := make(map[int64]*Job)
	jq := NewJobQueue(w, m)
	jqsrv := NewServer(jq)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs/enqueue", jqsrv.Enqueue)
	mux.HandleFunc("POST /jobs/dequeue", jqsrv.Dequeue)
	mux.HandleFunc("PATCH /jobs/{jobId}/conclude", jqsrv.Conclude)
	mux.HandleFunc("GET /jobs/{jobId}", jqsrv.Info)

	srv := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
	}

	log.Printf("Starting server on %s", address)
	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Unexpected HTTP server error: %v", err)
		}
	}
}
