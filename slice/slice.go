package slice

import (
	"sync"

	"github.com/rinx/golang-jobqueue/worker"
)

type Queue struct {
	sl []worker.JobFunc
	mu sync.Mutex
}

func NewQueue() worker.Queue {
	sl := make([]worker.JobFunc, 0, 0)

	return &Queue{
		sl: sl,
		mu: sync.Mutex{},
	}
}

func (q *Queue) Push(f worker.JobFunc) {
	q.mu.Lock()
	q.sl = append(q.sl, f)
	q.mu.Unlock()
}

func (q *Queue) Pop() (f worker.JobFunc) {
	q.mu.Lock()
	if len(q.sl) > 0 {
		f = q.sl[0]
		q.sl = q.sl[1:]
	}
	q.mu.Unlock()

	return f
}
