package discover

import (
	"errors"
	"io/fs"
	"os"
	"sync/atomic"
)

func NewFSExplorer(folder string) Explorer {
	e := &fsExplorer{
		fired:   false,
		count:   0,
		visited: make(chan string, 4),
		errors:  make(chan error),
		folder:  folder,
		done:    make(chan struct{}),
	}

	return e
}

type fsExplorer struct {
	fired bool

	count   uint64
	visited chan string
	errors  chan error

	done chan struct{}

	folder string
}

var ErrYetFired = errors.New("The explorer has yet been fired")

func (e *fsExplorer) Explore() error {
	if e.fired {
		return ErrYetFired
	}

	go e.explore()

	return nil
}

func (e *fsExplorer) explore() {
	defer e.dispose()
	fsys := os.DirFS(e.folder)
	fs.WalkDir(fsys, ".", e.dirWalker)

}

func (e *fsExplorer) dispose() {
	close(e.visited)
	close(e.errors)
	close(e.done)
}

func (e *fsExplorer) dirWalker(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if p == e.folder {
		return nil
	}

	if !d.IsDir() {
		e.visit(p)
	}

	return nil
}

func (e *fsExplorer) Errors() <-chan error {
	return e.errors
}

func (e *fsExplorer) Sleep() error {
	return nil
}

func (e *fsExplorer) Wake() error {
	return nil
}

func (e *fsExplorer) Visited() <-chan string {
	return e.visited
}

func (e *fsExplorer) visit(p string) {
	e.visited <- p
	atomic.AddUint64(&e.count, 1)
}

func (e *fsExplorer) VisitedNum() uint64 {
	return e.count
}

func (e *fsExplorer) Done() <-chan struct{} {
	return e.done
}
