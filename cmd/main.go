package main

import (
	"log"
	"os"
	"time"

	"github.com/filariow/polo/pkg/discover"
	"github.com/filariow/polo/pkg/execute"
	"github.com/schollz/progressbar/v3"
)

const folder = "/tmp/aa"

func getFolder() string {
	if len(os.Args) == 2 {
		return os.Args[1]
	}
	return folder
}

func main() {
	f := getFolder()

	ce := discover.NewFSExplorer(f)
	go countElements(ce)

	e := discover.NewFSExplorer(f)
	if err := e.Explore(); err != nil {
		log.Fatal(err)
	}

	c := execute.NewClamdExecutor(f, e.Visited(), 1)
	c.Execute()

	b := progressbar.NewOptions64(
		int64(ce.VisitedNum()),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetDescription("Scanned files..."),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(50*time.Millisecond),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetVisibility(true),
	)

	func() {
		for {
			select {
			case <-c.Done():
				showProgress(b, ce, c)
				return
			case <-time.After(200 * time.Millisecond):
				showProgress(b, ce, c)
			}
		}
	}()

	<-e.Done()
	<-ce.Done()

	time.Sleep(200 * time.Millisecond)
}

func countElements(e discover.Explorer) uint64 {
	err := e.Explore()
	if err != nil {
		log.Fatal(err)
	}

	func() {
		for {
			select {
			case <-e.Done():
				return
			case <-e.Visited():
			}
		}
	}()

	return e.VisitedNum()
}

func showProgress(bar *progressbar.ProgressBar, ex discover.Explorer, ec execute.Executor) {
	r := ex.VisitedNum()
	t := ec.Terminated()

	bar.ChangeMax64(int64(r))
	bar.Set64(int64(t))
}
