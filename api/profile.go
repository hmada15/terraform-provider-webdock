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

type (
	Profile struct {
		Slug  string `json:"slug"`
		Name  string `json:"name"`
		RAM   int    `json:"ram"`
		Disk  int    `json:"disk"`
		CPU   CPU    `json:"cpu"`
		Price Price  `json:"price"`
	}
	CPU struct {
		Cores   int `json:"cores"`
		Threads int `json:"threads"`
	}
	Price struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	}
)

func (c *Client) ListProfiles(ctx context.Context, locationId string) ([]Profile, error) {
	uri := BASE_URL + "profiles?locationId=" + locationId

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return []Profile{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return []Profile{}, errors.New("unexpected http error code received for geting Profile data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var profiles []Profile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return []Profile{}, err
	}

	return profiles, nil
}
