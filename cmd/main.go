package main

import (
	"log"

	"github.com/filariow/polo/pkg/discover"
)

func main() {
	e := discover.NewFSExplorer(".")

	err := e.Explore()
	if err != nil {
		log.Fatal(err)
	}

	for v := range e.Visited() {
		log.Println("visited", v)
	}

	log.Println("Completed")
}
