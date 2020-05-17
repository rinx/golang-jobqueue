package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/rinx/golang-jobqueue/worker"
)

type queue struct {
	inCh  chan worker.JobFunc
	outCh chan worker.JobFunc
}

func newQueue(ctx context.Context) worker.Queue {
	inCh := make(chan worker.JobFunc)
	outCh := make(chan worker.JobFunc)
	sl := make([]worker.JobFunc, 0, 0)

	go func() {
		defer close(inCh)
		defer close(outCh)

		for {
			if len(sl) > 0 {
				select {
				case <-ctx.Done():
					return
				case outCh <- sl[0]:
					sl = sl[1:]
				case in := <-inCh:
					sl = append(sl, in)
				}
			} else {
				select {
				case <-ctx.Done():
					return
				case in := <-inCh:
					sl = append(sl, in)
				}
			}
		}
	}()

	return &queue{
		inCh:  inCh,
		outCh: outCh,
	}
}

func (q *queue) Push(f worker.JobFunc) {
	q.inCh <- f
}

func (q *queue) Pop() worker.JobFunc {
	return <-q.outCh
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	n := 1000000

	q := newQueue(ctx)

	w := worker.NewWorker(ctx, q)
	w.Start()

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		func(i int) {
			q.Push(func() {
				defer wg.Done()
				fmt.Println(fmt.Sprintf("job %d", i))
			})
		}(i)
	}

	wg.Wait()

	cancel()
}
