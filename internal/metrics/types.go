package metrics

import "time"

type SystemMetrics struct {
	Memory    MemoryInfo      `json:"memory"`
	Disks     []DiskInfo      `json:"disks"`
	Uptime    time.Duration   `json:"uptime"`
	Kernel    string          `json:"kernel"`
	Services  []ServiceInfo   `json:"services"`
	Containers []ContainerInfo `json:"containers"`
	Timestamp time.Time       `json:"timestamp"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Active bool   `json:"active"`
}

type ContainerInfo struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Image   string   `json:"image"`
	Ports   []string `json:"ports"`
	Created int64    `json:"created"`
}
