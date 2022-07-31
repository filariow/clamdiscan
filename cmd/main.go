package main

import (
	"fmt"
	"log"
	"time"

	"github.com/filariow/polo/pkg/discover"
	"github.com/filariow/polo/pkg/execute"
)

func main() {
	e := discover.NewFSExplorer(".")

	err := e.Explore()
	if err != nil {
		log.Fatal(err)
	}

	v := e.Visited()
	c := execute.NewLogExecutor(v)
	c.Execute()

	for {
		time.Sleep(1 * time.Second)
		select {
		case <-c.Done():
			showProgress(e, c)
			fmt.Println("")
			return
		default:
			showProgress(e, c)
		}
	}

}

func showProgress(ex discover.Explorer, ec execute.Executor) {
	r := ex.VisitedNum()
	o := ec.ExecutedNum()

	fmt.Printf("visited/executed %d/%d\r", r, o)
}
