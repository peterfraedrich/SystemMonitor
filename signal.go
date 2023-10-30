package main

import (
	"os"
	"os/signal"
)

func CatchSignals(e chan Event) {

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	for {
		s := <-c
		e <- Event{
			Source: "OS",
			Type:   "SIGNAL",
			Event:  s,
		}
	}
}
