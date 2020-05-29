package ch

import (
	"github.com/rinx/golang-jobqueue/worker"
)

type Queue struct {
	ch chan worker.JobFunc
}

func NewQueue(buf int) worker.Queue {
	ch := make(chan worker.JobFunc, buf)

	return &Queue{
		ch: ch,
	}
}

func (q *Queue) Push(f worker.JobFunc) {
	q.ch <- f
}

func (q *Queue) Pop() worker.JobFunc {
	return <-q.ch
}
