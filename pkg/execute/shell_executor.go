package execute

import (
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"sync/atomic"
)

type shellExecutor struct {
	folder string
	source <-chan string

	started     uint64
	terminated  uint64
	done        chan struct{}
	concurrency int

	workers chan struct{}
	found   chan string
}

func NewShellExecutor(folder string, source <-chan string, t int) Executor {
	e := &shellExecutor{
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

func (l *shellExecutor) Execute() error {
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

func (l *shellExecutor) execute(e string) error {
	atomic.AddUint64(&l.started, 1)
	defer atomic.AddUint64(&l.terminated, 1)

	a := path.Join(l.folder, e)
	cmd := exec.Command("clamdscan", "--infected", "--no-summary", "--fdpass", a)
	cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {

		ec := cmd.ProcessState.ExitCode()
		switch ec {
		case 1:
			return nil
		case 2:
			// error parsing file
			return nil
		default:
			return err
		}
	}

	return nil
}

func (l *shellExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *shellExecutor) Started() uint64 {
	return l.started
}
func (l *shellExecutor) Terminated() uint64 {
	return l.terminated
}
