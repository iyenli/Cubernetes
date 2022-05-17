package cache

import (
	"runtime"

	dockertypes "github.com/docker/docker/api/types"
)

type CpuStatsCache interface {
	CalculateCpuPercent(containerID string, newCpuStats dockertypes.CPUStats) float64
}

func NewCpuStatsCache() CpuStatsCache {
	return &cpuStatsCache{
		cpuCount: runtime.NumCPU(),
		cache:    make(map[string]cpuStats),
	}
}

type cpuStatsCache struct {
	cpuCount int
	cache    map[string]cpuStats
}

type cpuStats struct {
	TotalUsage  uint64
	SystemUsage uint64
}

// calculate CpuPercent using newly update CpuStats and old stats in cache,
// then update stats in cache, return 0.0 if containerID not present.
func (c *cpuStatsCache) CalculateCpuPercent(containerID string, newCpuStats dockertypes.CPUStats) float64 {
	cpuPercent := 0.0

	oldStats, ok := c.cache[containerID]
	if ok {
		cpuDelta := float64(newCpuStats.CPUUsage.TotalUsage) - float64(oldStats.TotalUsage)
		systemDelta := float64(newCpuStats.SystemUsage) - float64(oldStats.SystemUsage)
		if cpuDelta > 0.0 && systemDelta > 0.0 {
			cpuPercent = (cpuDelta / systemDelta) * float64(c.cpuCount) * 100.0
		}
	}

	c.cache[containerID] = cpuStats{
		TotalUsage:  newCpuStats.CPUUsage.TotalUsage,
		SystemUsage: newCpuStats.SystemUsage,
	}

	return cpuPercent
}
