package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/filariow/polo/pkg/discover"
	"github.com/filariow/polo/pkg/execute"
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

	c := execute.NewShellExecutor(f, e.Visited(), 8)
	c.Execute()

	func() {
		for {
			select {
			case <-c.Done():
				showProgress(e, c)
				fmt.Println("")
				return
			case <-time.After(200 * time.Millisecond):
				showProgress(ce, c)
			}
		}
	}()

	<-c.Done()
	<-e.Done()

	<-ce.Done()
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

func showProgress(ex discover.Explorer, ec execute.Executor) {
	r := ex.VisitedNum()
	o := ec.ExecutedNum()

	fmt.Printf("executed/visited %d/%d\r", o, r)
}
