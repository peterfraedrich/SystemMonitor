package main

import (
	"sync"

	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInformation struct {
	Hostname        string
	BootTime        uint64
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
	HostID          string
	CPUVendorID     string
	CPUFamily       string
	CPUModel        string
	CPUID           string
	CPUCores        int32
	CPUMHZ          float64
	CPUCacheSize    int32
}

type SystemMetricsBasic struct {
	System struct {
		Uptime     uint64
		CPUVoltage uint16
	} `gorm:"embedded;embeddedPrefix:system_"`
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
		sysinfo.Hostname = stat.Hostname
		sysinfo.BootTime = stat.BootTime
		sysinfo.OS = stat.OS
		sysinfo.Platform = stat.Platform
		sysinfo.PlatformFamily = stat.PlatformFamily
		sysinfo.PlatformVersion = stat.PlatformVersion
		sysinfo.KernelVersion = stat.KernelVersion
		sysinfo.KernelArch = stat.KernelArch
		sysinfo.HostID = stat.HostID
		if err != nil {
			errs = append(errs, err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := cpu.Info()
		if err != nil {
			errs = append(errs, err)
			return
		}
		cpu := stat[0]
		sysinfo.CPUVendorID = cpu.VendorID
		sysinfo.CPUFamily = cpu.Family
		sysinfo.CPUModel = cpu.Model
		sysinfo.CPUID = cpu.CoreID
		sysinfo.CPUCores = cpu.Cores
		sysinfo.CPUMHZ = cpu.Mhz
		sysinfo.CPUCacheSize = cpu.CacheSize
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
		if len(stat) > 0 {
			copier.Copy(&metrics.CPU, &stat[0])
		}
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
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	stat, err := host.SensorsTemperatures()
	//	if err != nil {
	//		errs = append(errs, err)
	//		return
	//	}
	//	if len(stat) > 0 {
	//		copier.Copy(&metrics.Temps, &stat[0])
	//	}
	//}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		stat, err := host.Info()
		if err != nil {
			errs = append(errs, err)
			return
		}
		metrics.System.Uptime = stat.Uptime
		if err != nil {
			errs = append(errs, err)
		}
	}()
	wg.Wait()
	return metrics, errs
}
