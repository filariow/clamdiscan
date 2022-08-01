package execute

import (
	"log"
	"sync/atomic"
)

type logExecutor struct {
	source <-chan string

	started    uint64
	terminated uint64

	done chan struct{}
}

func NewLogExecutor(source <-chan string) Executor {
	return &logExecutor{
		source:     source,
		started:    0,
		terminated: 0,
		done:       make(chan struct{}),
	}
}

func (l *logExecutor) Execute() error {
	go func() {
		for e := range l.source {
			l.execute(e)
		}

		close(l.done)
	}()

	return nil
}

func (l *logExecutor) execute(e string) {
	atomic.AddUint64(&l.started, 1)
	defer atomic.AddUint64(&l.terminated, 1)

	log.Println(e)
}

func (l *logExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *logExecutor) Started() uint64 {
	return l.started
}
func (l *logExecutor) Terminated() uint64 {
	return l.terminated
}
