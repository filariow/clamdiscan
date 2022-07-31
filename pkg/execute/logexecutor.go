package execute

import (
	"sync/atomic"
)

type logExecutor struct {
	source <-chan string

	count uint64
	done  chan struct{}
}

func NewLogExecutor(source <-chan string) Executor {
	return &logExecutor{
		source: source,
		count:  0,
		done:   make(chan struct{}),
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
	atomic.AddUint64(&l.count, 1)
}

func (l *logExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *logExecutor) ExecutedNum() uint64 {
	return l.count
}
