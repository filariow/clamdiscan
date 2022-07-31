package main

import (
	"log"

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

	<-c.Done()
}
