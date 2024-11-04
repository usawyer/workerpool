package workerpool

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	workers  atomic.Int64
	jobChan  chan string
	quitChan chan struct{}
	wg       sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		jobChan:  make(chan string),
		quitChan: make(chan struct{}),
	}
}

func (wp *WorkerPool) StartWorker(id int) {
	for {
		select {
		case job, ok := <-wp.jobChan:
			if !ok {
				return
			}
			slog.Info(fmt.Sprintf("Worker %d processing: %s", id, job))
			time.Sleep(time.Second)
		case <-wp.quitChan:
			slog.Info(fmt.Sprintf("Worker %d removed", id))
			return
		}
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.wg.Add(1)
	id := wp.workers.Add(1)
	go func() {
		defer wp.wg.Done()
		wp.StartWorker(int(id))
	}()
	slog.Info(fmt.Sprintf("Worker %d created", id))
}

func (wp *WorkerPool) AddJob(job string) {
	wp.jobChan <- job
}

func (wp *WorkerPool) RemoveWorker() {
	if wp.workers.Load() == 0 {
		slog.Info("No workers to remove")
		return
	}
	wp.workers.Add(-1)
	wp.quitChan <- struct{}{}
}

func (wp *WorkerPool) StopAll() {
	close(wp.quitChan)
	wp.WaitAll()
	wp.CloseJobChannel()
}

func (wp *WorkerPool) WaitAll() {
	wp.wg.Wait()
}

func (wp *WorkerPool) CloseJobChannel() {
	close(wp.jobChan)
}

func (wp *WorkerPool) GetWorkersNum() int64 {
	return wp.workers.Load()
}

func (wp *WorkerPool) GetJob() string {
	return <-wp.jobChan
}
