package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupServer() *Server {
	// wire up dependencies
	w := make([]int64, 0)
	m := make(map[int64]*Job)
	jq := NewJobQueue(w, m)
	jqsrv := NewServer(jq)
	return jqsrv
}

// todo: need more/better tests
func TestEnqueue(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus int
		expectedBody   Job
	}{
		{"blue sky", `{"Type":"TIME_CRITICAL"}`, http.StatusCreated, Job{Type: TimeCritical, ID: 1, Status: Queued}},
		{"bad input", `{"Typo":"TIME_CRITICAL"}`, http.StatusBadRequest, Job{}},
	}

	jqsrv := setupServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/jobs/enqueue", bytes.NewReader([]byte(tt.input)))
			w := httptest.NewRecorder()
			jqsrv.Enqueue(w, req)

			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expect %d, got %d", tt.expectedStatus, status)
			}

			result := Job{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			if tt.expectedBody != result {
				t.Errorf("expected %v, got %v", tt.expectedBody, result)
			}
		})
	}
}

// todo: need more/better tests
func TestDequeue(t *testing.T) {
	tests := []struct {
		name           string
		consumer       string
		expectedStatus int
		expectedBody   Job
	}{
		{"blue sky", "123", http.StatusOK, Job{Type: TimeCritical, ID: 1, Status: InProgress, QueueConsumer: 123}},
		{"bad queue consumer ID", "abc", http.StatusBadRequest, Job{}},
	}

	w := []int64{1}
	m := map[int64]*Job{1: {Type: TimeCritical, ID: 1, Status: Queued}}
	jq := NewJobQueue(w, m)
	jqsrv := NewServer(jq)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/jobs/dequeue", nil)
			req.Header.Add(QueueConsumerHdr, tt.consumer)
			w := httptest.NewRecorder()
			jqsrv.Dequeue(w, req)

			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expect %d, got %d", tt.expectedStatus, status)
			}

			result := Job{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			if tt.expectedBody != result {
				t.Errorf("expected %v, got %v", tt.expectedBody, result)
			}
		})
	}
}

// todo: finish test
// func TestConclude(t *testing.T) {
//
// }

// todo: finish test
// func TestInfo(t *testing.T) {
//
// }
