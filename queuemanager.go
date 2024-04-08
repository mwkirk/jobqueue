package main

import (
	"errors"
	"fmt"
	"sync"
)

type JobStatus int

const (
	Unspecified JobStatus = iota
	Queued
	InProgress
	Concluded
)

var TimeCritical = "TIME_CRITICAL"
var NotTimeCritical = "NOT_TIME_CRITICAL"

var ErrorWaitQueueEmpty = errors.New("wait queue empty")

type Job struct {
	ID            int64     `json:"id"`
	Type          string    `json:"type"`
	Status        JobStatus `json:"status"`
	QueueConsumer int64     `json:"queue_consumer,omitempty"`
}

type JobQueue struct {
	l      sync.RWMutex
	jobNum int64
	wait   []int64
	jobs   map[int64]*Job
}

func NewJobQueue(q []int64, m map[int64]*Job) *JobQueue {
	return &JobQueue{
		wait: q,
		jobs: m,
	}
}

func (jq *JobQueue) Enqueue(j *Job) error {
	// Bare minimum validation
	if j.Type != TimeCritical && j.Type != NotTimeCritical {
		return fmt.Errorf("invalid job type: %s", j.Type)
	}

	jq.l.Lock()
	defer jq.l.Unlock()
	jq.jobNum++
	j.ID = jq.jobNum
	j.Status = Queued
	jq.wait = append(jq.wait, j.ID)
	jq.jobs[j.ID] = j
	return nil
}

func (jq *JobQueue) Dequeue(n int64) (*Job, error) {
	jq.l.Lock()
	defer jq.l.Unlock()
	if len(jq.wait) == 0 {
		return nil, ErrorWaitQueueEmpty
	}

	// dequeue the job from the wait queue
	id := jq.wait[0]
	jq.wait = jq.wait[1:]

	j, ok := jq.jobs[id]
	if !ok {
		// It's still best in this case to remove the job id from the wait queue
		return nil, fmt.Errorf("job %d not found", id)
	}

	// update the job
	j.Status = InProgress
	j.QueueConsumer = n
	return j, nil
}

func (jq *JobQueue) Conclude(id int64) (*Job, error) {
	jq.l.Lock()
	defer jq.l.Unlock()

	j, ok := jq.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job %d not found when concluding job", id)
	}

	j.Status = Concluded
	return j, nil
}

func (jq *JobQueue) Info(id int64) (*Job, error) {
	jq.l.RLock()
	defer jq.l.RUnlock()

	j, ok := jq.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job %d not found during info lookup", id)
	}

	return j, nil
}
