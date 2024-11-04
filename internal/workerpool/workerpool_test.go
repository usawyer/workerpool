package workerpool_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/usawyer/workerpool/internal/workerpool"
	"testing"
)

func TestWorkerPool_AddWorker(t *testing.T) {
	wp := workerpool.NewWorkerPool()
	wp.AddWorker()
	assert.Equal(t, int64(1), wp.GetWorkersNum())
}

func TestWorkerPool_AddJob(t *testing.T) {
	wp := workerpool.NewWorkerPool()
	job := "test job"
	go wp.AddJob(job)
	assert.Equal(t, job, wp.GetJob())
}

func TestWorkerPool_RemoveWorker(t *testing.T) {
	wp := workerpool.NewWorkerPool()
	wp.AddWorker()
	assert.Equal(t, int64(1), wp.GetWorkersNum())
	wp.RemoveWorker()
	assert.Equal(t, int64(0), wp.GetWorkersNum())
}

func TestWorkerPool_StopAll(t *testing.T) {
	wp := workerpool.NewWorkerPool()
	wp.AddWorker()
	assert.Equal(t, int64(1), wp.GetWorkersNum())
	wp.AddJob("test job")
	wp.StopAll()
	assert.Panics(t, func() { wp.AddJob("another job") })
}
