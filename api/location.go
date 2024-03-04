package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/hmada15/terraform-provider-webdock/helper"
)

type Location struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (c *Client) ListLocations(ctx context.Context) ([]Location, error) {
	uri := BASE_URL + "locations"

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return []Location{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return []Location{}, errors.New("unexpected http error code received for geting Location data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var locations []Location
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return []Location{}, err
	}

	return locations, nil
}
