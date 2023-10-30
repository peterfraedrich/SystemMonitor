package main

import (
	"sync"

	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInformation struct {
	Host struct {
		Hostname             string `json:"hostname"`
		Uptime               uint64 `json:"uptime"`
		BootTime             uint64 `json:"bootTime"`
		Procs                uint64 `json:"procs"`           // number of processes
		OS                   string `json:"os"`              // ex: freebsd, linux
		Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
		PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
		PlatformVersion      string `json:"platformVersion"` // version of the complete OS
		KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
		KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
		VirtualizationSystem string `json:"virtualizationSystem"`
		VirtualizationRole   string `json:"virtualizationRole"` // guest or host
		HostID               string `json:"hostId"`             // ex: uuid
	} `gorm:"embedded;embeddedPrefix:host_"`
	CPU struct {
		CPU        int32    `json:"cpu"`
		VendorID   string   `json:"vendorId"`
		Family     string   `json:"family"`
		Model      string   `json:"model"`
		Stepping   int32    `json:"stepping"`
		PhysicalID string   `json:"physicalId"`
		CoreID     string   `json:"coreId"`
		Cores      int32    `json:"cores"`
		ModelName  string   `json:"modelName"`
		Mhz        float64  `json:"mhz"`
		CacheSize  int32    `json:"cacheSize"`
		Flags      []string `json:"flags"`
		Microcode  string   `json:"microcode"`
	} `gorm:"embedded;embeddedPrefix:cpu_"`
}

type SystemMetricsBasic struct {
	CPU struct {
		CPU       string  `json:"cpu"`
		User      float64 `json:"user"`
		System    float64 `json:"system"`
		Idle      float64 `json:"idle"`
		Nice      float64 `json:"nice"`
		Iowait    float64 `json:"iowait"`
		Irq       float64 `json:"irq"`
		Softirq   float64 `json:"softirq"`
		Steal     float64 `json:"steal"`
		Guest     float64 `json:"guest"`
		GuestNice float64 `json:"guestNice"`
	} `gorm:"embedded;embeddedPrefix:cpu_"`
	Memory struct {
		Total          uint64  `json:"total"`
		Available      uint64  `json:"available"`
		Used           uint64  `json:"used"`
		UsedPercent    float64 `json:"usedPercent"`
		Free           uint64  `json:"free"`
		Active         uint64  `json:"active"`
		Inactive       uint64  `json:"inactive"`
		Wired          uint64  `json:"wired"`
		Laundry        uint64  `json:"laundry"`
		Buffers        uint64  `json:"buffers"`
		Cached         uint64  `json:"cached"`
		WriteBack      uint64  `json:"writeBack"`
		Dirty          uint64  `json:"dirty"`
		WriteBackTmp   uint64  `json:"writeBackTmp"`
		Shared         uint64  `json:"shared"`
		Slab           uint64  `json:"slab"`
		Sreclaimable   uint64  `json:"sreclaimable"`
		Sunreclaim     uint64  `json:"sunreclaim"`
		PageTables     uint64  `json:"pageTables"`
		SwapCached     uint64  `json:"swapCached"`
		CommitLimit    uint64  `json:"commitLimit"`
		CommittedAS    uint64  `json:"committedAS"`
		HighTotal      uint64  `json:"highTotal"`
		HighFree       uint64  `json:"highFree"`
		LowTotal       uint64  `json:"lowTotal"`
		LowFree        uint64  `json:"lowFree"`
		SwapTotal      uint64  `json:"swapTotal"`
		SwapFree       uint64  `json:"swapFree"`
		Mapped         uint64  `json:"mapped"`
		VmallocTotal   uint64  `json:"vmallocTotal"`
		VmallocUsed    uint64  `json:"vmallocUsed"`
		VmallocChunk   uint64  `json:"vmallocChunk"`
		HugePagesTotal uint64  `json:"hugePagesTotal"`
		HugePagesFree  uint64  `json:"hugePagesFree"`
		HugePagesRsvd  uint64  `json:"hugePagesRsvd"`
		HugePagesSurp  uint64  `json:"hugePagesSurp"`
		HugePageSize   uint64  `json:"hugePageSize"`
		AnonHugePages  uint64  `json:"anonHugePages"`
	} `gorm:"embedded;embeddedPrefix:mem_"`
	Temps struct {
		SensorKey   string  `json:"sensorKey"`
		Temperature float64 `json:"temperature"`
		High        float64 `json:"sensorHigh"`
		Critical    float64 `json:"sensorCritical"`
	} `gorm:"embedded;embeddedPrefix:temp_"`
}

func GetSystemInformation() (SystemInformation, []error) {
	sysinfo := SystemInformation{}
	errs := []error{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := host.Info()
		if err != nil {
			errs = append(errs, err)
			return
		}
		err = copier.Copy(&sysinfo.Host, &stat)
		if err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		stat, err := cpu.Info()
		if err != nil {
			errs = append(errs, err)
			return
		}
		copier.Copy(&sysinfo.CPU, &stat[0])
		if err != nil {
			errs = append(errs, err)
		}
	}()
	wg.Wait()
	if len(errs) != 0 {
		return sysinfo, errs
	}
	return sysinfo, nil
}

func GetSystemMetricsBasic() (SystemMetricsBasic, []error) {
	metrics := SystemMetricsBasic{}
	errs := []error{}
	wg := sync.WaitGroup{}
	// CPU
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := cpu.Times(false)
		if err != nil {
			errs = append(errs, err)
			return
		}
		copier.Copy(&metrics.CPU, &stat[0])
	}()
	// MEM
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := mem.VirtualMemory()
		if err != nil {
			errs = append(errs, err)
			return
		}
		copier.Copy(&metrics.Memory, &stat)
	}()
	// TEMPS
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := host.SensorsTemperatures()
		if err != nil {
			errs = append(errs, err)
			return
		}
		copier.Copy(&metrics.Temps, &stat[0])
	}()
	wg.Wait()
	return metrics, errs
}
