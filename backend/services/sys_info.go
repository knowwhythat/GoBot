package services

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type LSysInfo struct {
	MemAll         uint64
	MemFree        uint64
	MemUsed        uint64
	MemUsedPercent float64

	DiskAll         uint64
	DiskFree        uint64
	DiskUsed        uint64
	DiskUsedPercent float64

	CpuUsedPercent   float64
	CpuPhysicalCount int
	CpuLogicalCount  int
	OS               string
	OSPlatform       string
	Arch             string
	Hostname         string
	HostId           string
}

func GetSysInfo() (info LSysInfo) {
	unit := uint64(1024 * 1024) // MB

	v, _ := mem.VirtualMemory()

	info.MemAll = v.Total
	info.MemFree = v.Free
	info.MemUsed = info.MemAll - info.MemFree
	// 注：使用SwapMemory或VirtualMemory，在不同系统中使用率不一样，因此直接计算一次
	info.MemUsedPercent = float64(info.MemUsed) / float64(info.MemAll) * 100.0
	info.MemAll /= unit
	info.MemUsed /= unit
	info.MemFree /= unit

	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH

	d, _ := disk.Usage("/")
	n, _ := host.Info()
	info.CpuPhysicalCount, _ = cpu.Counts(false)
	info.CpuLogicalCount, _ = cpu.Counts(true)
	info.DiskAll = d.Total / unit / 1024
	info.DiskFree = d.Free / unit / 1024
	info.DiskUsed = d.Used / unit / 1024
	info.DiskUsedPercent = d.UsedPercent
	info.Hostname = n.Hostname
	info.OSPlatform = n.Platform + n.PlatformFamily + n.PlatformVersion + n.KernelArch
	info.HostId = n.HostID
	// 获取200ms内的CPU信息，太短不准确，也可以获几秒内的，但这样会有延时，因为要等待
	cc, _ := cpu.Percent(time.Second*1, false)
	info.CpuUsedPercent = cc[0]

	return
}
