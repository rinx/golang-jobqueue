package worker

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)

type JobFunc func()

type Queue interface {
	Push(f JobFunc)
	Pop() JobFunc
}

type Worker interface {
	Start()
}

type worker struct {
	ctx   context.Context
	queue Queue
}

func NewWorker(ctx context.Context, q Queue) Worker {
	return &worker{
		ctx:   ctx,
		queue: q,
	}
}

func (w *worker) Start() {
	wg := sync.WaitGroup{}

	sm := semaphore.NewWeighted(1)

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				wg.Wait()
				return
			default:
			}

			f := w.queue.Pop()
			if f != nil {
				wg.Add(1)

				err := sm.Acquire(w.ctx, 1)
				if err != nil {
					fmt.Println(err)
				}

				go func(f JobFunc) {
					defer wg.Done()
					defer sm.Release(1)

					f()
				}(f)
			}
		}
	}()
}
