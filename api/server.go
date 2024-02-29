package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hmada15/terraform-provider-webdock/helper"

	"net/http"
)

type Server struct {
	Slug                   string `json:"slug"`
	Name                   string `json:"name"`
	Date                   string `json:"date"`
	Location               string `json:"location"`
	Image                  string `json:"image"`
	Profile                string `json:"profile"`
	Ipv4                   string `json:"ipv4,omitempty"`
	Ipv6                   string `json:"ipv6,omitempty"`
	Status                 string `json:"status"`
	Virtualization         string `json:"virtualization"`
	WebServer              string `json:"webServer"`
	SnapshotRunTime        int64  `json:"snapshotRunTime"`
	Description            string `json:"description"`
	WordPressLockDown      bool   `json:"WordPressLockDown"`
	SSHPasswordAuthEnabled bool   `json:"SSHPasswordAuthEnabled"`
	NextActionDate         string `json:"nextActionDate"`
}

func (c *Client) GetServerBYSlug(ctx context.Context, slug string) (Server, error) {
	uri := BASE_URL + "servers/" + slug

	resp, err := helper.NewWebdockRequest(http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Server{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Server{}, errors.New("unexpected http error code received for geting server")
	}
	defer resp.Body.Close()

	var servers Server
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return Server{}, err
	}

	return servers, nil
}
