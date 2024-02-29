package api

const BASE_URL = "https://api.webdock.io/v1/"

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}
