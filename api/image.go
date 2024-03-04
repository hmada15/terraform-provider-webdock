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
	Image struct {
		Slug       string `json:"slug"`
		Name       string `json:"name"`
		WebServer  string `json:"webServer"`
		PhpVersion string `json:"phpVersion"`
	}
)

func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	uri := BASE_URL + "images"

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return []Image{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return []Image{}, errors.New("unexpected http error code received for geting serverImages data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var images []Image
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return []Image{}, err
	}

	return images, nil
}
