package instrumentation

import (
	"runtime"
)

type varzMemoryStats struct {
	BytesAllocatedHeap  uint64 `json:"numBytesAllocatedHeap"`
	BytesAllocatedStack uint64 `json:"numBytesAllocatedStack"`
	BytesAllocated      uint64 `json:"numBytesAllocated"`
	NumMallocs          uint64 `json:"numMallocs"`
	NumFrees            uint64 `json:"numFrees"`
	LastGCPauseTimeNS   uint64 `json:"lastGCPauseTimeNS"`
}

func mapMemStats(stats *runtime.MemStats) varzMemoryStats {
	return varzMemoryStats{BytesAllocatedHeap: stats.HeapAlloc,
		BytesAllocatedStack: stats.StackInuse,
		BytesAllocated:      stats.Alloc,
		NumMallocs:          stats.Mallocs,
		NumFrees:            stats.Frees,
		LastGCPauseTimeNS:   stats.PauseNs[(stats.NumGC+255)%256]}
}

type VarzMessage struct {
	Name          string            `json:"name"`
	NumCpus       int               `json:"numCPUS"`
	NumGoRoutines int               `json:"numGoRoutines"`
	MemoryStats   varzMemoryStats   `json:"memoryStats"`
	Tags          map[string]string `json:"tags"`
	Contexts      []Context         `json:"contexts"`
}

func NewVarzMessage(name string, ipAddress string, instrumentables []Instrumentable) (*VarzMessage, error) {
	contexts := make([]Context, len(instrumentables))
	for i, instrumentable := range instrumentables {
		contexts[i] = instrumentable.Emit()
	}
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	tags := map[string]string{
		"ip": ipAddress,
	}

	return &VarzMessage{name, runtime.NumCPU(), runtime.NumGoroutine(), mapMemStats(memStats), tags, contexts}, nil
}
