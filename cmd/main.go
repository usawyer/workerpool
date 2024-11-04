package main

import (
	"fmt"
	"github.com/usawyer/workerpool/internal/workerpool"
	"github.com/usawyer/workerpool/pkg/logger"
	"time"
)

const (
	loggerLevel = "info"
)

func main() {
	logger.New(loggerLevel)

	p := workerpool.NewWorkerPool()

	// Добавление ворекеров
	p.AddWorker()
	p.AddWorker()

	// Добавление заданий
	go func() {
		for i := 0; i < 10; i++ {
			p.AddJob(fmt.Sprintf("Job %d", i+1))
		}
	}()

	time.Sleep(2 * time.Second)

	// Удаление одного воркера
	p.RemoveWorker()

	//Добавление новых воркеров
	p.AddWorker()
	p.AddWorker()
	p.AddWorker()

	// Добавляем еще заданий
	go func() {
		for i := 10; i < 20; i++ {
			p.AddJob(fmt.Sprintf("Job %d", i+1))
		}
		p.CloseJobChannel()
	}()

	p.WaitAll()
}
