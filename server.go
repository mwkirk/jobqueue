package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

var QueueConsumerHdr = "X-QUEUE-CONSUMER"

type Server struct {
	mgr *JobQueue
}

func NewServer(jq *JobQueue) *Server {
	return &Server{
		mgr: jq,
	}
}

// Enqueue submits a job to the queue.
func (s *Server) Enqueue(w http.ResponseWriter, r *http.Request) {
	// todo: check for content type

	var j Job
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
		return
	}

	err = s.mgr.Enqueue(&j)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(j)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
	log.Printf("Enqueuing job %v\n", j)
}

// Dequeue takes the next job in the wait queue and returns it to the calling consumer
func (s *Server) Dequeue(w http.ResponseWriter, r *http.Request) {
	qcStr := r.Header.Get(QueueConsumerHdr)
	if len(qcStr) == 0 {
		http.Error(w, "Missing header "+QueueConsumerHdr, http.StatusBadRequest)
		return
	}

	qc, err := strconv.Atoi(qcStr)
	if err != nil {
		http.Error(w, "Invalid consumer ID value in header "+QueueConsumerHdr, http.StatusBadRequest)
		return
	}

	j, err := s.mgr.Dequeue(int64(qc))
	if err != nil {
		if errors.Is(err, ErrorWaitQueueEmpty) {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			http.Error(w, "Unable to dequeue job", http.StatusInternalServerError)
			return
		}
	}

	out, err := json.Marshal(j)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
	log.Printf("Dequeued job %v\n", j)
}

// Conclude marks the job with the provided job ID as done. The job ID is sent as a
// path parameter, e.g. "/jobs/{jobId}/conclude".
func (s *Server) Conclude(w http.ResponseWriter, r *http.Request) {
	jobIdStr := r.PathValue("jobId")
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job ID value in path", http.StatusBadRequest)
		return
	}

	qcStr := r.Header.Get(QueueConsumerHdr)
	if len(qcStr) == 0 {
		http.Error(w, "Missing header "+QueueConsumerHdr, http.StatusBadRequest)
		return
	}

	// todo: see below
	_, err = strconv.Atoi(qcStr)
	if err != nil {
		http.Error(w, "Invalid consumer ID value in header "+QueueConsumerHdr, http.StatusBadRequest)
		return
	}

	// todo: This job could already be concluded, but we allow it to be concluded again
	// todo: Should we check if the queue consumer concluding is the same one that dequeued?
	j, err := s.mgr.Conclude(int64(jobId))
	if err != nil {
		http.Error(w, "Job ID does not exist", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(j)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
	log.Printf("Concluded job %v\n", j)
}

// Info returns info about a job given the job ID provided via path parameter, i.e. "/jobs/{jobId}".
func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	jobIdStr := r.PathValue("jobId")
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job ID value in path", http.StatusBadRequest)
		return
	}

	j, err := s.mgr.Info(int64(jobId))
	if err != nil {
		http.Error(w, "Job ID does not exist", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(j)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched info for job %v\n", j)
}
