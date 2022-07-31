package execute

import (
	"time"
)

type logExecutor struct {
	source <-chan string
	count  int64
	done   chan struct{}
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
	// fmt.Println(e)
	l.count++

	time.Sleep(40 * time.Millisecond)
}

func (l *logExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *logExecutor) ExecutedNum() int64 {
	return l.count
}
