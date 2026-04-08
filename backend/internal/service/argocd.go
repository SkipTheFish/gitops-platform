package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ArgoCDService struct {
	server string
	token  string
	client *http.Client
}

func NewArgoCDService(server, token string) *ArgoCDService {
	return &ArgoCDService{
		server: strings.TrimRight(server, "/"),
		token:  token,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type ArgoCDAppStatus struct {
	SyncStatus   string
	HealthStatus string
	Revision     string
}

func (s *ArgoCDService) GetApplicationStatus(appName string) (*ArgoCDAppStatus, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", s.server, appName)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get argocd app status failed: %s", resp.Status)
	}

	var raw struct {
		Status struct {
			Sync struct {
				Status   string `json:"status"`
				Revision string `json:"revision"`
			} `json:"sync"`
			Health struct {
				Status string `json:"status"`
			} `json:"health"`
		} `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &ArgoCDAppStatus{
		SyncStatus:   raw.Status.Sync.Status,
		HealthStatus: raw.Status.Health.Status,
		Revision:     raw.Status.Sync.Revision,
	}, nil
}

func (s *ArgoCDService) SyncApplication(appName string) error {
	url := fmt.Sprintf("%s/api/v1/applications/%s/sync", s.server, appName)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(`{}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("sync argocd app failed: %s", resp.Status)
	}
	return nil
}
