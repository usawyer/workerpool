package workerpool

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool
// workers - количество воркеров
// jobChan - канал с джобами
// quitChan - канал для завершения работы воркера
// wg - примитив синхронизации
type WorkerPool struct {
	workers  atomic.Int64
	jobChan  chan string
	quitChan chan struct{}
	wg       sync.WaitGroup
}

// NewWorkerPool создание нового ВоркерПула
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		jobChan:  make(chan string),
		quitChan: make(chan struct{}),
	}
}

// StartWorker запускает работу воркера
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

// AddWorker добавляет нового воркера
func (wp *WorkerPool) AddWorker() {
	wp.wg.Add(1)
	id := wp.workers.Add(1)
	go func() {
		defer wp.wg.Done()
		wp.StartWorker(int(id))
	}()
	slog.Info(fmt.Sprintf("Worker %d created", id))
}

// AddJob добавляет новую джобу
func (wp *WorkerPool) AddJob(job string) {
	wp.jobChan <- job
}

// RemoveWorker удаляет воркера из пула
func (wp *WorkerPool) RemoveWorker() {
	if wp.workers.Load() == 0 {
		slog.Info("No workers to remove")
		return
	}
	wp.workers.Add(-1)
	wp.quitChan <- struct{}{}
}

// StopAll останавливает работу всех воркеров
func (wp *WorkerPool) StopAll() {
	close(wp.quitChan)
	wp.WaitAll()
	wp.CloseJobChannel()
}

// WaitAll дожидается окончания выполнения джобов воркерами
func (wp *WorkerPool) WaitAll() {
	wp.wg.Wait()
}

// CloseJobChannel закрывает канал джобов
func (wp *WorkerPool) CloseJobChannel() {
	close(wp.jobChan)
}

// GetWorkersNum возвращает текущее количество воркеров
func (wp *WorkerPool) GetWorkersNum() int64 {
	return wp.workers.Load()
}

// GetJob возвращает текущий джоб
func (wp *WorkerPool) GetJob() string {
	return <-wp.jobChan
}
