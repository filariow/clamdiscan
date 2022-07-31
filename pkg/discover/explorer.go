package discover

type Explorer interface {
	Explore() error
	Sleep() error
	Wake() error

	Done() <-chan struct{}

	VisitedNum() uint64
	Visited() <-chan string
	Errors() <-chan error
}
