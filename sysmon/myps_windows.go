//go:build windows
// +build windows

package sysmon

import (
	"syscall"
	"unsafe"
)

var (
	kernel32, _                 = syscall.LoadLibrary("kernel32.dll")
	procGetSystemTimes, _       = syscall.GetProcAddress(kernel32, "GetSystemTimes")
	procGlobalMemoryStatusEx, _ = syscall.GetProcAddress(kernel32, "GlobalMemoryStatusEx")
	procGetLastError, _         = syscall.GetProcAddress(kernel32, "GetLastError")
)

func GetRAMPercent() (float64, error) {
	var memInfo struct {
		cbSize                  uint32
		dwMemoryLoad            uint32
		ullTotalPhys            uint64
		ullAvailPhys            uint64
		ullTotalPageFile        uint64
		ullAvailPageFile        uint64
		ullTotalVirtual         uint64
		ullAvailVirtual         uint64
		ullAvailExtendedVirtual uint64
	}
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := syscall.SyscallN(procGlobalMemoryStatusEx, uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		if r0, _, _ := syscall.SyscallN(procGetLastError); r0 != 0 {
			return 0, syscall.Errno(r0)
		}
		return 0, nil
	}
	return float64(memInfo.dwMemoryLoad), nil
}

func cpuTimes() ([]timesStat, error) {
	var lpIdleTime, lpKernelTime, lpUserTime struct {
		DwLowDateTime  uint32
		DwHighDateTime uint32
	}

	r, _, _ := syscall.SyscallN(
		procGetSystemTimes,
		uintptr(unsafe.Pointer(&lpIdleTime)),
		uintptr(unsafe.Pointer(&lpKernelTime)),
		uintptr(unsafe.Pointer(&lpUserTime)),
	)
	if r == 0 {
		if r0, _, _ := syscall.SyscallN(procGetLastError); r0 != 0 {
			return nil, syscall.Errno(r0)
		}
		return nil, nil
	}

	LOT := float64(0.0000001)
	HIT := (LOT * 4294967296.0)
	idle := ((HIT * float64(lpIdleTime.DwHighDateTime)) + (LOT * float64(lpIdleTime.DwLowDateTime)))
	user := ((HIT * float64(lpUserTime.DwHighDateTime)) + (LOT * float64(lpUserTime.DwLowDateTime)))
	kernel := ((HIT * float64(lpKernelTime.DwHighDateTime)) + (LOT * float64(lpKernelTime.DwLowDateTime)))
	system := (kernel - idle)

	ret := []timesStat{
		{
			CPU:    "cpu-total",
			Idle:   float64(idle),
			User:   float64(user),
			System: float64(system),
		},
	}

	return ret, nil
}
