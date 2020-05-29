package main

import (
	"context"
	"testing"
	"time"

	"github.com/rinx/golang-jobqueue/ch"
	"github.com/rinx/golang-jobqueue/csc"
	"github.com/rinx/golang-jobqueue/slice"
	"github.com/rinx/golang-jobqueue/worker"
)

const (
	bufn = 10

	inBufn  = 10
	outBufn = 10
	dur     = 1 * time.Millisecond

	nworker = 1
)

var (
	j = func() {
		time.Sleep(100 * time.Millisecond)
	}
)

func BenchmarkQueue_Channel(b *testing.B) {
	ctx, _ := context.WithCancel(context.Background())

	q := ch.NewQueue(bufn)

	for i := 0; i < nworker; i++ {
		w := worker.NewWorker(ctx, q)
		w.Start()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(j)
	}
}

func BenchmarkQueue_Slice(b *testing.B) {
	ctx, _ := context.WithCancel(context.Background())

	q := slice.NewQueue()

	for i := 0; i < nworker; i++ {
		w := worker.NewWorker(ctx, q)
		w.Start()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(j)
	}
}

func BenchmarkQueue_ChannelSliceChannel(b *testing.B) {
	ctx, _ := context.WithCancel(context.Background())

	q := csc.NewQueue(ctx, inBufn, outBufn, dur)

	for i := 0; i < nworker; i++ {
		w := worker.NewWorker(ctx, q)
		w.Start()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(j)
	}
}
