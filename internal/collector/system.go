package collector

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/dima/gentoo-monitor/internal/metrics"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemCollector struct{}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (c *SystemCollector) GetMemoryInfo() (metrics.MemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return metrics.MemoryInfo{}, err
	}

	return metrics.MemoryInfo{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}, nil
}

func (c *SystemCollector) GetDiskInfo() ([]metrics.DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []metrics.DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, metrics.DiskInfo{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	return disks, nil
}

func (c *SystemCollector) GetUptime() (time.Duration, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return 0, err
	}
	return time.Duration(uptime) * time.Second, nil
}

func (c *SystemCollector) GetKernelVersion() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.KernelVersion, nil
}

func (c *SystemCollector) GetSystemdServices() ([]metrics.ServiceInfo, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-pager", "--no-legend")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var services []metrics.ServiceInfo
	lines := strings.Split(out.String(), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ".service")
		status := fields[2]
		active := status == "active"

		services = append(services, metrics.ServiceInfo{
			Name:   name,
			Status: status,
			Active: active,
		})
	}

	return services, nil
}
