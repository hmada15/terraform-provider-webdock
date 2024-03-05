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

type PublicKey struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	Created string `json:"created"`
}
type PublicKeyRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func (c *Client) GetPublicKeyById(ctx context.Context, id string) (PublicKey, error) {
	uri := BASE_URL + "account/publicKeys"

	resp, err := helper.NewWebdockRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return PublicKey{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return PublicKey{}, errors.New("unexpected http error code received for geting PublicKey data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var publicKeys []PublicKey
	if err := json.NewDecoder(resp.Body).Decode(&publicKeys); err != nil {
		return PublicKey{}, err
	}

	publicKey := PublicKey{}
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return PublicKey{}, err
	}
	for _, key := range publicKeys {
		if key.ID == idStr {
			publicKey = key
			break
		}
	}

	return publicKey, nil
}

func (c *Client) CreatePublicKey(ctx context.Context, publicKeyRequest PublicKeyRequest) (PublicKey, error) {
	uri := BASE_URL + "account/publicKeys"

	jsonPayload, err := json.Marshal(publicKeyRequest)
	if err != nil {
		return PublicKey{}, err
	}

	resp, err := helper.NewWebdockRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return PublicKey{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return PublicKey{}, errors.New("unexpected http error code received for creating publickey status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var publicKey PublicKey
	if err := json.NewDecoder(resp.Body).Decode(&publicKey); err != nil {
		return PublicKey{}, err
	}

	return publicKey, nil
}

func (c *Client) DeletePublicKey(ctx context.Context, id string) error {
	uri := BASE_URL + "account/publicKeys/" + id

	resp, err := helper.NewWebdockRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("unexpected http error code received for deleting publickey status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	return nil
}
