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

	count       uint64
	done        chan struct{}
	concurrency int

	workers chan struct{}
}

func NewShellExecutor(folder string, source <-chan string, t int) Executor {
	e := &shellExecutor{
		folder:      folder,
		source:      source,
		count:       0,
		done:        make(chan struct{}),
		concurrency: t,
		workers:     make(chan struct{}, t),
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
	atomic.AddUint64(&l.count, 1)

	a := path.Join(l.folder, e)
	cmd := exec.Command("clamscan", "--infected", "--no-summary", "--bell", a)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		if ec := cmd.ProcessState.ExitCode(); ec == 1 {
			return nil
		}
		return err
	}

	return nil
}
func (l *shellExecutor) Done() <-chan struct{} {
	return l.done
}

func (l *shellExecutor) ExecutedNum() uint64 {
	return l.count
}
