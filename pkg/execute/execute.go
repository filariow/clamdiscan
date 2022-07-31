package execute

type Executor interface {
	Execute() error

	Started() uint64
	Terminated() uint64

	Done() <-chan struct{}
}
