package execute

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/baruwa-enterprise/clamd"
)

var ErrInvalidClamdAddress = fmt.Errorf("Invalid clamd address")

type clamdExecutor struct {
	folder string
	source <-chan string

	started     uint64
	terminated  uint64
	done        chan struct{}
	concurrency int

	workers chan struct{}
	found   chan string

	clamdAddress string
	clamdClient  *clamd.Client
}

func NewClamdExecutor(address string, folder string, source <-chan string, t int) Executor {
	e := &clamdExecutor{
		folder:       folder,
		source:       source,
		started:      0,
		terminated:   0,
		done:         make(chan struct{}),
		concurrency:  t,
		workers:      make(chan struct{}, t),
		found:        make(chan string),
		clamdAddress: address,
	}
	return e
}

func (l *clamdExecutor) parseAddress() (string, string, error) {
	ss := strings.Split(l.clamdAddress, "://")
	if len(ss) != 2 {
		return "", "", ErrInvalidClamdAddress
	}

	return ss[0], ss[1], nil
}

func (l *clamdExecutor) Execute() error {
	h, a, err := l.parseAddress()
	if err != nil {
		return fmt.Errorf("error parsing address '%v': %w", l.clamdAddress, err)
	}

	c, err := clamd.NewClient(h, a)
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
