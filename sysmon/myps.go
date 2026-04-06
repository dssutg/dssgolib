package sysmon

import (
	"github.com/dssutg/dssgolib/utils"
	"time"
)

type timesStat struct {
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
}

// GetCPUPercent samples CPU times and returns per-core usage.
func GetCPUPercent() ([]float64, error) {
	t1, err := cpuTimes()
	if err != nil {
		return nil, err
	}

	time.Sleep(300 * time.Millisecond)

	t2, err := cpuTimes()
	if err != nil {
		return nil, err
	}

	// If counts differ, only compute for the minimum length to stay safe.
	n := min(len(t2), len(t1))
	percent := make([]float64, n)

	for i := range n {
		oldStats := t1[i]
		newStats := t2[i]

		oldIdle := oldStats.Idle + oldStats.Iowait
		newIdle := newStats.Idle + newStats.Iowait

		oldTotal := oldStats.User + oldStats.Nice + oldStats.System + oldStats.Idle + oldStats.Iowait + oldStats.Irq + oldStats.Softirq + oldStats.Steal + oldStats.Guest + oldStats.GuestNice
		newTotal := newStats.User + newStats.Nice + newStats.System + newStats.Idle + newStats.Iowait + newStats.Irq + newStats.Softirq + newStats.Steal + newStats.Guest + newStats.GuestNice

		totalDelta := newTotal - oldTotal
		idleDelta := newIdle - oldIdle

		if totalDelta <= 0 {
			percent[i] = 0
			continue
		}
		percent[i] = utils.Clamp((totalDelta-idleDelta)/totalDelta*100.0, 0, 100)
	}

	return percent, nil
}
