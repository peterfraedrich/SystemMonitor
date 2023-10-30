package main

import (
	"sync"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessMetricsAggregate struct {
	ProcCount      uint64 `json:"procCount"`
	ThreadCount    uint64 `json:"threadCount"`
	ProcForeground uint64 `json:"procForeground"`
	ProcBackground uint64 `json:"procBackground"`
	ProcRunning    uint64 `json:"procRunning"`
	ProcSleeping   uint64 `json:"procSleeping"`
	ProcStopped    uint64 `json:"procStopped"`
	ProcIdle       uint64 `json:"procIdle"`
	ProcZombie     uint64 `json:"procZombie"`
	ProcWaiting    uint64 `json:"procWaiting"`
	ProcLocked     uint64 `json:"procLocked"`
	OpenFiles      uint64 `json:"openFiles"`
}

func GetAggregateProcessMetrics() (ProcessMetricsAggregate, []error) {
	metrics := ProcessMetricsAggregate{}
	errs := []error{}
	wg := sync.WaitGroup{}
	proclist, err := process.Processes()
	if err != nil {
		errs = append(errs, err)
	}
	// proc count
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.ProcCount = uint64(len(proclist))
	}()
	// thread count
	wg.Add(1)
	go func() {
		defer wg.Done()
		var c uint64 = 0
		for _, p := range proclist {
			t, _ := p.NumThreads()
			c = uint64(t) + c
		}
	}()
	// proc state
	wg.Add(1)
	go func() {
		defer wg.Done()
		var fg uint64
		var bg uint64
		var run uint64
		var sleep uint64
		var stop uint64
		var idle uint64
		var zombie uint64
		var wait uint64
		var lock uint64
		for _, p := range proclist {
			t, err := p.Background()
			if err != nil {
				continue
			}
			if t {
				bg++
			} else {
				fg++
			}
			status, err := p.Status()
			if err != nil {
				continue
			}
			switch status[0] {
			case "R":
				run++
			case "S":
				sleep++
			case "T":
				stop++
			case "I":
				idle++
			case "Z":
				zombie++
			case "W":
				wait++
			case "L":
				lock++
			}
		}
		metrics.ProcBackground = bg
		metrics.ProcForeground = fg
		metrics.ProcIdle = idle
		metrics.ProcSleeping = sleep
		metrics.ProcStopped = stop
		metrics.ProcIdle = idle
		metrics.ProcZombie = zombie
		metrics.ProcWaiting = wait
		metrics.ProcLocked = lock
	}()
	// Open Files
	wg.Add(1)
	go func() {
		defer wg.Done()
		var fds uint64 = 0
		for _, p := range proclist {
			openFiles, _ := p.OpenFiles()
			if openFiles != nil {
				fds += uint64(len(openFiles))
			}
		}
		metrics.OpenFiles = fds
	}()
	wg.Wait()
	return metrics, errs
}
