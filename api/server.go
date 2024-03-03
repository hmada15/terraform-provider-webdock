package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

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
	WordPressLockDown      bool   `json:"WordPressLockDown"`
	SSHPasswordAuthEnabled bool   `json:"SSHPasswordAuthEnabled"`
}

type ServerRequest struct {
	Name           string `json:"name"`
	Slug           string `json:"slug,omitempty"`
	LocationID     string `json:"locationId"`
	ProfileSlug    string `json:"profileSlug"`
	Virtualization string `json:"virtualization,omitempty"`
	ImageSlug      string `json:"imageSlug"`
}

func (c *Client) GetServerBYSlug(ctx context.Context, slug string) (Server, error) {
	uri := BASE_URL + "servers/" + slug

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Server{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Server{}, errors.New("unexpected http error code received for geting server data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var server Server
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return Server{}, err
	}

	return server, nil
}

func (c *Client) CreateServer(ctx context.Context, serverRequest ServerRequest) (Server, error) {
	uri := BASE_URL + "servers"

	jsonPayload, err := json.Marshal(serverRequest)
	if err != nil {
		return Server{}, err
	}

	resp, err := helper.NewWebdockRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return Server{}, err
	}
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return Server{}, errors.New("unexpected http error code received for creating server status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var server Server
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return Server{}, err
	}

	return server, nil
}

func (c *Client) UpdateServer(ctx context.Context, slug string, serverRequest ServerRequest) (Server, error) {
	uri := BASE_URL + "servers/" + slug

	jsonPayload, err := json.Marshal(serverRequest)
	if err != nil {
		return Server{}, err
	}

	resp, err := helper.NewWebdockRequest(ctx, http.MethodPatch, uri, jsonPayload, c.token)
	if err != nil {
		return Server{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Server{}, errors.New("unexpected http error code received for updating server status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var server Server
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return Server{}, err
	}

	return server, nil
}

func (c *Client) DeleteServer(ctx context.Context, slug string) error {
	uri := BASE_URL + "servers/" + slug

	resp, err := helper.NewWebdockRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("unexpected http error code received for deleting server status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) ServerExist(ctx context.Context, slug string) (string, error) {
	uri := BASE_URL + "servers/" + slug

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return helper.NO, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("unexpected http error code received for checking if server exist status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	return helper.YES, nil
}
