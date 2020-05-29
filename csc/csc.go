package csc

import (
	"context"
	"time"

	"github.com/rinx/golang-jobqueue/worker"
)

type Queue struct {
	inCh  chan worker.JobFunc
	outCh chan worker.JobFunc
}

func NewQueue(ctx context.Context, inBuf, outBuf int, dur time.Duration) worker.Queue {
	inCh := make(chan worker.JobFunc, inBuf)
	outCh := make(chan worker.JobFunc, outBuf)
	sl := make([]worker.JobFunc, 0, 0)

	ticker := time.NewTicker(dur)

	go func() {
		defer close(inCh)
		defer close(outCh)

		for {
			if len(sl) > 0 {
				select {
				case <-ctx.Done():
					return
				case in := <-inCh:
					sl = append(sl, in)
				case outCh <- sl[0]:
					sl = sl[1:]
				}
			}

			select {
			case <-ctx.Done():
				return
			case in := <-inCh:
				sl = append(sl, in)
			case <-ticker.C:
			}
		}
	}()

	return &Queue{
		inCh:  inCh,
		outCh: outCh,
	}
}

func (q *Queue) Push(f worker.JobFunc) {
	q.inCh <- f
}

func (q *Queue) Pop() worker.JobFunc {
	return <-q.outCh
}
