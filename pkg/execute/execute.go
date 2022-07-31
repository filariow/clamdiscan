package execute

type Executor interface {
	Execute() error

	ExecutedNum() int64

	Done() <-chan struct{}
}
