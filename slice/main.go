package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/rinx/golang-jobqueue/worker"
)

type queue struct {
	sl []worker.JobFunc
	mu sync.Mutex
}

func newQueue() worker.Queue {
	sl := make([]worker.JobFunc, 0)

	return &queue{
		sl: sl,
		mu: sync.Mutex{},
	}
}

func (q *queue) Push(f worker.JobFunc) {
	q.mu.Lock()
	q.sl = append(q.sl, f)
	q.mu.Unlock()
}

func (q *queue) Pop() (f worker.JobFunc) {
	q.mu.Lock()
	if len(q.sl) > 0 {
		f = q.sl[0]
		q.sl = q.sl[1:]
	}
	q.mu.Unlock()

	return f
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	n := 1000000

	q := newQueue()

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
