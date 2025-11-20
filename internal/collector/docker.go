package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/dima/gentoo-monitor/internal/metrics"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerCollector struct {
	client *client.Client
}

func NewDockerCollector() (*DockerCollector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerCollector{client: cli}, nil
}

func (c *DockerCollector) GetContainers() ([]metrics.ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containers, err := c.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []metrics.ContainerInfo
	for _, container := range containers {
		ports := c.formatPorts(container.Ports)
		name := container.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		result = append(result, metrics.ContainerInfo{
			ID:      container.ID[:12],
			Name:    name,
			State:   container.State,
			Status:  container.Status,
			Image:   container.Image,
			Ports:   ports,
			Created: container.Created,
		})
	}

	return result, nil
}

func (c *DockerCollector) formatPorts(ports []types.Port) []string {
	var result []string
	for _, port := range ports {
		if port.PublicPort > 0 {
			result = append(result, fmt.Sprintf("%d:%d/%s", port.PublicPort, port.PrivatePort, port.Type))
		} else {
			result = append(result, fmt.Sprintf("%d/%s", port.PrivatePort, port.Type))
		}
	}
	return result
}

func (c *DockerCollector) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
