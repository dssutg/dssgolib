//go:build linux
// +build linux

package sysmon

import (
	"bufio"
	"github.com/dssutg/dssgolib/utils"
	"io"
	"os"
	"strings"
)

const clocksPerSec = 100

func cpuTimes() ([]timesStat, error) {
	filename := "/proc/stat"
	lines := readLines(filename, 1)
	ret := make([]timesStat, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) < 8 {
			continue
		}

		if !strings.HasPrefix(fields[0], "cpu") {
			continue
		}

		cpu := fields[0]
		if cpu == "cpu" {
			cpu = "cpu-total"
		}
		user, err := utils.ParseFloat64(fields[1])
		if err != nil {
			continue
		}
		nice, err := utils.ParseFloat64(fields[2])
		if err != nil {
			continue
		}
		system, err := utils.ParseFloat64(fields[3])
		if err != nil {
			continue
		}
		idle, err := utils.ParseFloat64(fields[4])
		if err != nil {
			continue
		}
		iowait, err := utils.ParseFloat64(fields[5])
		if err != nil {
			continue
		}
		irq, err := utils.ParseFloat64(fields[6])
		if err != nil {
			continue
		}
		softirq, err := utils.ParseFloat64(fields[7])
		if err != nil {
			continue
		}

		ct := &timesStat{
			CPU:     cpu,
			User:    user / clocksPerSec,
			Nice:    nice / clocksPerSec,
			System:  system / clocksPerSec,
			Idle:    idle / clocksPerSec,
			Iowait:  iowait / clocksPerSec,
			Irq:     irq / clocksPerSec,
			Softirq: softirq / clocksPerSec,
		}
		if len(fields) > 8 { // Linux >= 2.6.11
			steal, err := utils.ParseFloat64(fields[8])
			if err != nil {
				continue
			}
			ct.Steal = steal / clocksPerSec
		}
		if len(fields) > 9 { // Linux >= 2.6.24
			guest, err := utils.ParseFloat64(fields[9])
			if err != nil {
				continue
			}
			ct.Guest = guest / clocksPerSec
		}
		if len(fields) > 10 { // Linux >= 3.2.0
			guestNice, err := utils.ParseFloat64(fields[10])
			if err != nil {
				continue
			}
			ct.GuestNice = guestNice / clocksPerSec
		}

		ret = append(ret, *ct)
	}

	return ret, nil
}

func GetRAMPercent() (float64, error) {
	var total, free, buffers, cached, sReclaimable uint64
	lines := readLines("/proc/meminfo", -1)

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.ReplaceAll(value, " kB", "")
		switch key {
		case "MemTotal":
			t, err := utils.ParseDecimalUint64(value)
			if err != nil {
				return 0, err
			}
			total = t * 1024
		case "MemFree":
			t, err := utils.ParseDecimalUint64(value)
			if err != nil {
				return 0, err
			}
			free = t * 1024
		case "Buffers":
			t, err := utils.ParseDecimalUint64(value)
			if err != nil {
				return 0, err
			}
			buffers = t * 1024
		case "Cached":
			t, err := utils.ParseDecimalUint64(value)
			if err != nil {
				return 0, err
			}
			cached = t * 1024
		case "SReclaimable":
			t, err := utils.ParseDecimalUint64(value)
			if err != nil {
				return 0, err
			}
			sReclaimable = t * 1024
		}
	}

	used := total - free - buffers - cached - sReclaimable
	usedPercent := float64(used) / float64(total) * 100.0

	return usedPercent, nil
}

func readLines(filename string, n int) []string {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				ret = append(ret, strings.Trim(line, "\n"))
			}
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret
}
