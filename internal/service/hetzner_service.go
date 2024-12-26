package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// HetznerService handles operations related to Hetzner.
type HetznerService struct {
	client *hcloud.Client
}

// NewHetznerService initializes the HetznerService with the client.
func NewHetznerService(token string) (*HetznerService, error) {
	if token == "" {
		return nil, errors.New("HCLOUD_TOKEN is required")
	}

	client := hcloud.NewClient(hcloud.WithToken(token))
	if client == nil {
		return nil, fmt.Errorf("hetzner client not initialized")
	}

	return &HetznerService{client: client}, nil
}

// GetAllServers retrieves all servers from Hetzner.
func (s *HetznerService) GetAllServers() ([]*hcloud.Server, error) {

	ctx := context.Background()
	servers, err := s.client.Server.All(ctx)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// Get one server by id
func (s *HetznerService) GetServerByID(id int64) (*hcloud.Server, error) {
	ctx := context.Background()
	server, _, err := s.client.Server.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// CreateServer creates a new server on Hetzner.
func (s *HetznerService) CreateServer(name string, image string, serverType string) (*hcloud.Server, error) {
	ctx := context.Background()
	result, _, err := s.client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name: name,
		ServerType: &hcloud.ServerType{
			Name: serverType,
		},
		Image: &hcloud.Image{
			Name: image,
		},
	})
	if err != nil {
		return nil, err
	}
	err = s.client.Action.WaitFor(ctx, result.Action)
	if err != nil {
		return nil, err
	}
	return result.Server, nil
}

// Update server name and lables
func (s *HetznerService) UpdateServer(id int64, name string, labels map[string]string) (*hcloud.Server, error) {
	ctx := context.Background()
	server, _, err := s.client.Server.Update(ctx, &hcloud.Server{ID: id}, hcloud.ServerUpdateOpts{
		Name:   name,
		Labels: labels,
	})
	if err != nil {
		return nil, err
	}

	return server, nil

}

func (s *HetznerService) GetServerMetricsByID(id int64, start, end string) (*hcloud.ServerMetrics, error) {
	ctx := context.Background()

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	metrics, _, err := s.client.Server.GetMetrics(ctx, &hcloud.Server{ID: id}, hcloud.ServerGetMetricsOpts{
		Start: startTime,
		End:   endTime,
		Step:  60,
		Types: []hcloud.ServerMetricType{
			hcloud.ServerMetricCPU,
			hcloud.ServerMetricNetwork,
			hcloud.ServerMetricDisk,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %v", err)
	}

	return metrics, nil
}

// Delete server
func (s *HetznerService) DeleteServer(id int64) error {
	ctx := context.Background()
	result, _, err := s.client.Server.DeleteWithResult(ctx, &hcloud.Server{ID: id})

	if err != nil {
		return err
	}
	err = s.client.Action.WaitFor(ctx, result.Action)
	return err
}

// LogOperation logs operations to the database or log file.
func (s *HetznerService) LogOperation(action string, resource interface{}) {
	// Replace with actual database logic.
	log.Printf("Action: %s, Resource: %+v, Timestamp: %s\n", action, resource, time.Now().Format(time.RFC3339))
}
