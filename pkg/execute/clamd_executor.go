package execute

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"github.com/baruwa-enterprise/clamd"
)

type clamdExecutor struct {
	folder string
	source <-chan string

	started     uint64
	terminated  uint64
	done        chan struct{}
	concurrency int

	workers     chan struct{}
	found       chan string
	clamdClient *clamd.Client
}

func NewClamdExecutor(folder string, source <-chan string, t int) Executor {
	e := &clamdExecutor{
		folder:      folder,
		source:      source,
		started:     0,
		terminated:  0,
		done:        make(chan struct{}),
		concurrency: t,
		workers:     make(chan struct{}, t),
		found:       make(chan string),
	}
	return e
}

func (l *clamdExecutor) Execute() error {
	c, err := clamd.NewClient("unix", "/var/run/clamav/clamd.ctl")
	if err != nil {
		return err
	}
	l.clamdClient = c

	go func() {
		var wg sync.WaitGroup

		for e := range l.source {
			l.workers <- struct{}{}
			go func(h string) {
				defer wg.Done()
				wg.Add(1)

				if err := l.execute(h); err != nil {
					log.Println(err)
				}
				<-l.workers
			}(e)
		}

		wg.Wait()
		close(l.done)
	}()
	return nil
}

func (l *clamdExecutor) execute(e string) error {
	atomic.AddUint64(&l.started, 1)
	defer atomic.AddUint64(&l.terminated, 1)
	a := path.Join(l.folder, e)

	rr, err := l.clamdClient.Fildes(context.TODO(), a)
	if err != nil {
		return err
	}

	for _, r := range rr {
		if r.Status != "OK" {
			fmt.Fprintf(os.Stderr, "File: %v, Status: %v\n", a, r.Status)
		}
	}

	return nil
}

func (l *clamdExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *clamdExecutor) Started() uint64 {
	return l.started
}
func (l *clamdExecutor) Terminated() uint64 {
	return l.terminated
}
