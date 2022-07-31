package discover

import (
	"errors"
	"io/fs"
	"os"
)

func NewFSExplorer(folder string) Explorer {
	e := &fsExplorer{
		fired:   false,
		toVisit: []string{folder},
		count:   0,
		visited: make(chan string),
		errors:  make(chan error),
		folder:  folder,
	}

	return e
}

type fsExplorer struct {
	fired   bool
	toVisit []string

	count   int64
	visited chan string
	errors  chan error

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
	for {
		if len(e.toVisit) == 0 {
			break
		}

		p := e.toVisit[0]
		e.toVisit = e.toVisit[1:]

		fsys := os.DirFS(p)
		fs.WalkDir(fsys, ".", e.dirWalker)
	}

	close(e.visited)
	close(e.errors)
}

func (e *fsExplorer) dirWalker(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if p == e.folder {
		return nil
	}

	if d.IsDir() {
		e.toVisit = append(e.toVisit, p)
		return nil
	}

	e.visit(p)
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
	e.count++
}

func (e *fsExplorer) VisitedNum() int64 {
	return e.count
}
