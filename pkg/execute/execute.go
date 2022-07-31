package execute

import "fmt"

type Executor interface {
	Execute() error

	Done() <-chan struct{}
}

type logExecutor struct {
	source <-chan string
	done   chan struct{}
}

func NewLogExecutor(source <-chan string) Executor {
	return &logExecutor{
		source: source,
		done:   make(chan struct{}),
	}
}

func (l *logExecutor) Execute() error {
	go l.execute()
	return nil
}

func (l *logExecutor) execute() {
	for e := range l.source {
		fmt.Println(e)
	}

	close(l.done)

}

func (l *logExecutor) Done() <-chan struct{} {
	return l.done
}
