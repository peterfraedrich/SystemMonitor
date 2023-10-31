package main

import (
	"os"
	"os/signal"
	"syscall"
)

var signals = map[os.Signal]int{
	syscall.SIGINT:  2,
	syscall.SIGHUP:  1,
	syscall.SIGSEGV: 11,
}

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
		for k, v := range signals {
			if k == s {
				os.Exit(v)
			}
		}
	}
}
