package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/rinx/golang-jobqueue/worker"
)

type queue struct {
	ch chan worker.JobFunc
}

func newQueue(buf int) worker.Queue {
	ch := make(chan worker.JobFunc, buf)

	return &queue{
		ch: ch,
	}
}

func (q *queue) Push(f worker.JobFunc) {
	q.ch <- f
}

func (q *queue) Pop() worker.JobFunc {
	return <-q.ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	n := 1000000
	bufn := 10

	q := newQueue(bufn)

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
