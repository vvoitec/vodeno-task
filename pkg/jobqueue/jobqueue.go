// Package jobqueue implements basic in-memory job queue for asynchronous processing.
package jobqueue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type (
	Worker struct {
		ctx      context.Context
		handlers sync.Map
		queue    chan payload
		logger   *log.Logger
	}

	Queue struct {
		queue chan payload
	}

	JobLabel string

	HandlerFunc func(context.Context, interface{}) error

	payload struct {
		data  interface{}
		label JobLabel
	}
)

func NewWorker(ctx context.Context, logger *log.Logger) *Worker {
	return &Worker{
		ctx:    ctx,
		queue:  make(chan payload, 1000),
		logger: logger,
	}
}

// AddHandler registers new job queues handler.
func (w *Worker) AddHandler(label JobLabel, handler HandlerFunc) {
	w.handlers.Store(label, handler)
}

func (w *Worker) GetQueue() *Queue {
	return &Queue{queue: w.queue}
}

// Run runs worker and consumes queued messages.
// TODO: implement retries with some backoff strategy
// TODO: implement backpressure handling
func (w *Worker) Run() {
	defer close(w.queue)
	for {
		select {
		case <-w.ctx.Done():
			return
		case job := <-w.queue:
			func() {
				handler, ok := w.handlers.Load(job.label)
				if !ok {
					panic(fmt.Sprintf("failed to match job label: %s to any handler", job.label))
				}
				ctx, cancel := context.WithTimeout(w.ctx, 10*time.Second)
				defer cancel()
				if err := handler.(HandlerFunc)(ctx, job.data); err != nil {
					w.logger.Printf("job: %s failed with: %s", job.label, err.Error())
				}
			}()
		}
	}
}

// Enqueue dispatched job to job queue.
func (q *Queue) Enqueue(label JobLabel, data interface{}) {
	q.queue <- payload{label: label, data: data}
}
