package execute

type Executor interface {
	Execute() error

	ExecutedNum() uint64

	Done() <-chan struct{}
}
