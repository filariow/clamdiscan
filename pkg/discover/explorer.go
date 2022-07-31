package discover

type Explorer interface {
	Explore() error
	Sleep() error
	Wake() error

	VisitedNum() int64
	Visited() <-chan string
	Errors() <-chan error
}
