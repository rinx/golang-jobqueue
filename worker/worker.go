package worker

import (
	"context"
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
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			default:
			}

			f := w.queue.Pop()
			if f != nil {
				f()
			}
		}
	}()
}
