package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	Source string
	Type   string
	Event  interface{}
}

type Service struct {
	db     *gorm.DB
	events chan Event
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Start() {
	s.events = make(chan Event, 64)
	go s.Handle()
	go CatchSignals(s.events)
	go s.Main()
}

func (s *Service) Main() {
	for {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			info, errs := GetSystemInformation()
			if len(errs) != 0 {
				for _, e := range errs {
					s.NewEvent("GetSystemInformation", "ERROR", e)
				}
				return
			}
			s.NewEvent("GetSystemInformation", "SystemInformation", info)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			metrics, errs := GetSystemMetricsBasic()
			if len(errs) != 0 {
				for _, e := range errs {
					s.NewEvent("GetSystemMetricsBasic", "ERROR", e)
				}
				return
			}
			s.NewEvent("GetSystemMetricsBasic", "SystemMetricsBasic", metrics)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			proc, errs := GetAggregateProcessMetrics()
			if len(errs) != 0 {
				for _, e := range errs {
					s.NewEvent("GetAggregateProcessMetrics", "ERROR", e)
				}
				return
			}
			s.NewEvent("GetAggregateProcessMetrics", "ProcessMetricsAggregate", proc)
		}()
		time.Sleep(time.Duration(CONFIG.Frequency) * time.Second)
	}
}

func (s *Service) Handle() {
	for {
		e := <-s.events
		if CONFIG.LogToStdout {
			fmt.Printf("%+v\n", e.Event)
		}
		switch e.Type {
		case "SystemInformation":
			s.db.Save(&SystemInfo{
				SystemInformation: e.Event.(SystemInformation),
			})
		case "SystemMetricsBasic":
			s.db.Save(&BasicMetrics{
				SystemMetricsBasic: e.Event.(SystemMetricsBasic),
			})
		case "ProcessMetricsAggregate":
			s.db.Save(&ProcessMetrics{
				ProcessMetricsAggregate: e.Event.(ProcessMetricsAggregate),
			})
		case "EVENT":
			s.db.Save(&EventsLog{
				Source:  e.Source,
				Type:    e.Type,
				Content: e.Event.(string),
			})
		case "ERROR":
			s.db.Save(&ErrorLog{
				Source: e.Source,
				Error:  e.Event.(error).Error(),
			})
		case "SIGNAL":
			s.db.Save(&EventsLog{
				Source:  e.Source,
				Type:    e.Type,
				Content: e.Event.(os.Signal).String(),
			})
		}
	}
}

func (s *Service) NewEvent(source string, t string, event interface{}) {
	s.events <- Event{
		Source: source,
		Type:   t,
		Event:  event,
	}
}
