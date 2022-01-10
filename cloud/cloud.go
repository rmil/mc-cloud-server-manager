package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type (
	Config struct {
		Token              string `json:"token"`
		AppName            string `json:"-"`
		AppVersion         string `json:"-"`
		selectedDatacentre int
	}

	Clouder struct {
		client *hcloud.Client
		conf   Config
		dc     Datacentres
	}

	Datacentres struct {
		Datacentres    []hcloud.Datacenter `json:"datacenters"`
		Recommendation int                 `json:"recommendation"`
	}
)

var ErrRecommendationNotFound = errors.New("could not find recommended datacentre")

// NewCloud creates a client to Hetzner with the API token.
func NewCloud(conf Config) (*Clouder, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.hetzner.cloud/v1/datacenters", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make new HTTP request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get datacentres: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	dc := Datacentres{}
	err = json.Unmarshal(body, &dc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &Clouder{
		client: hcloud.NewClient(hcloud.WithToken(conf.Token), hcloud.WithApplication(conf.AppName, conf.AppVersion)),
		conf:   conf,
		dc:     dc,
	}, nil
}
