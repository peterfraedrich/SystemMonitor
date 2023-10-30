package main

import "gorm.io/gorm"

type SystemInfo struct {
	gorm.Model
	SystemInformation `gorm:"embedded"`
}

type BasicMetrics struct {
	gorm.Model
	SystemMetricsBasic `gorm:"embedded"`
}

type ProcessMetrics struct {
	gorm.Model
	ProcessMetricsAggregate `gorm:"embedded"`
}

type EventsLog struct {
	gorm.Model
	Source  string
	Type    string
	Content string
}

type ErrorLog struct {
	gorm.Model
	Source string
	Error  string
}
