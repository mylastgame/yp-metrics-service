package collector

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
)

const (
	Alloc          = "Alloc"
	BuckHashSys    = "BuckHashSys"
	Frees          = "Frees"
	GCCPUFraction  = "GCCPUFraction"
	GCSys          = "GCSys"
	HeapAlloc      = "HeapAlloc"
	HeapIdle       = "HeapIdle"
	HeapInuse      = "HeapInuse"
	HeapObjects    = "HeapObjects"
	HeapReleased   = "HeapReleased"
	HeapSys        = "HeapSys"
	LastGC         = "LastGC"
	Lookups        = "Lookups"
	MCacheInuse    = "MCacheInuse"
	MCacheSys      = "MCacheSys"
	MSpanInuse     = "MSpanInuse"
	MSpanSys       = "MSpanSys"
	Mallocs        = "Mallocs"
	NextGC         = "NextGC"
	NumForcedGC    = "NumForcedGC"
	NumGC          = "NumGC"
	OtherSys       = "OtherSys"
	PauseTotalNs   = "PauseTotalNs"
	StackInuse     = "StackInuse"
	StackSys       = "StackSys"
	Sys            = "Sys"
	TotalAlloc     = "TotalAlloc"
	PollCount      = "PollCount"
	RandomValue    = "RandomValue"
	TotalMemory    = "TotalMemory"
	FreeMemory     = "FreeMemory"
	CPUutilization = "CPUutilization"
)

type Collector struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Collector {
	return &Collector{storage: storage}
}

func (c *Collector) Collect() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	c.storage.SaveGauge(Alloc, float64(m.Alloc))
	c.storage.SaveGauge(BuckHashSys, float64(m.BuckHashSys))
	c.storage.SaveGauge(Frees, float64(m.Frees))
	c.storage.SaveGauge(GCCPUFraction, float64(m.GCCPUFraction))
	c.storage.SaveGauge(GCSys, float64(m.GCSys))
	c.storage.SaveGauge(HeapAlloc, float64(m.HeapAlloc))
	c.storage.SaveGauge(HeapIdle, float64(m.HeapIdle))
	c.storage.SaveGauge(HeapInuse, float64(m.HeapInuse))
	c.storage.SaveGauge(HeapObjects, float64(m.HeapObjects))
	c.storage.SaveGauge(HeapReleased, float64(m.HeapReleased))
	c.storage.SaveGauge(HeapSys, float64(m.HeapSys))
	c.storage.SaveGauge(LastGC, float64(m.LastGC))
	c.storage.SaveGauge(Lookups, float64(m.Lookups))
	c.storage.SaveGauge(MSpanInuse, float64(m.MSpanInuse))
	c.storage.SaveGauge(MSpanSys, float64(m.MSpanSys))
	c.storage.SaveGauge(MCacheInuse, float64(m.MCacheInuse))
	c.storage.SaveGauge(MCacheSys, float64(m.MCacheSys))
	c.storage.SaveGauge(Mallocs, float64(m.Mallocs))
	c.storage.SaveGauge(NextGC, float64(m.NextGC))
	c.storage.SaveGauge(NumForcedGC, float64(m.NumForcedGC))
	c.storage.SaveGauge(NumGC, float64(m.NumGC))
	c.storage.SaveGauge(OtherSys, float64(m.OtherSys))
	c.storage.SaveGauge(PauseTotalNs, float64(m.PauseTotalNs))
	c.storage.SaveGauge(StackInuse, float64(m.StackInuse))
	c.storage.SaveGauge(StackSys, float64(m.StackSys))
	c.storage.SaveGauge(Sys, float64(m.Sys))
	c.storage.SaveGauge(TotalAlloc, float64(m.TotalAlloc))

	c.storage.SaveCounter(PollCount, 1)
	c.storage.SaveGauge(RandomValue, rand.Float64())
}

func (c *Collector) CollectGOPSUtil() error {
	m, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get virtual memory fail: %v", err)
	}

	c.storage.SaveGauge(TotalMemory, float64(m.Total))
	c.storage.SaveGauge(FreeMemory, float64(m.Free))

	cpus, err := cpu.Percent(0, true)
	if err != nil {
		return fmt.Errorf("get cpu percent fail: %v", err)
	}

	for i, cp := range cpus {
		c.storage.SaveGauge(fmt.Sprintf("%s%d", CPUutilization, i+1), cp)
	}

	return nil
}
