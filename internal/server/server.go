package server

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dima/gentoo-monitor/internal/collector"
	"github.com/dima/gentoo-monitor/internal/metrics"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	systemCollector *collector.SystemCollector
	dockerCollector *collector.DockerCollector
	mux             *http.ServeMux
}

func NewServer() (*Server, error) {
	dockerCollector, err := collector.NewDockerCollector()
	if err != nil {
		log.Printf("Warning: Docker collector unavailable: %v", err)
		dockerCollector = nil
	}

	s := &Server{
		systemCollector: collector.NewSystemCollector(),
		dockerCollector: dockerCollector,
		mux:             http.NewServeMux(),
	}

	s.routes()
	return s, nil
}

func (s *Server) routes() {
	s.mux.Handle("/", http.FileServer(http.FS(staticFiles)))
	s.mux.HandleFunc("/api/metrics", s.handleMetrics)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	systemMetrics := metrics.SystemMetrics{
		Timestamp: time.Now(),
	}

	memory, err := s.systemCollector.GetMemoryInfo()
	if err != nil {
		log.Printf("Error getting memory info: %v", err)
	} else {
		systemMetrics.Memory = memory
	}

	disks, err := s.systemCollector.GetDiskInfo()
	if err != nil {
		log.Printf("Error getting disk info: %v", err)
	} else {
		systemMetrics.Disks = disks
	}

	uptime, err := s.systemCollector.GetUptime()
	if err != nil {
		log.Printf("Error getting uptime: %v", err)
	} else {
		systemMetrics.Uptime = uptime
	}

	kernel, err := s.systemCollector.GetKernelVersion()
	if err != nil {
		log.Printf("Error getting kernel version: %v", err)
	} else {
		systemMetrics.Kernel = kernel
	}

	services, err := s.systemCollector.GetSystemdServices()
	if err != nil {
		log.Printf("Error getting systemd services: %v", err)
	} else {
		systemMetrics.Services = services
	}

	if s.dockerCollector != nil {
		containers, err := s.dockerCollector.GetContainers()
		if err != nil {
			log.Printf("Error getting containers: %v", err)
		} else {
			systemMetrics.Containers = containers
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(systemMetrics); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Close() error {
	if s.dockerCollector != nil {
		return s.dockerCollector.Close()
	}
	return nil
}
