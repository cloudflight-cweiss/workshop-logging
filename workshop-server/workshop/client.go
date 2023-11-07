package workshop

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client   http.Client
	endpoint string
}

func NewClient(endpoint string) *Client {
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}
	return &Client{
		client:   http.Client{Timeout: time.Duration(5) * time.Second},
		endpoint: endpoint,
	}
}

func (c *Client) getApiUrl(api string) string {
	return fmt.Sprintf("%s/api/%s", c.endpoint, api)
}

func (c *Client) doGet(api string) error {
	resp, err := c.client.Get(c.getApiUrl(api))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) doPost(api string) error {
	_, err := c.client.Post(c.getApiUrl(api), "", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Login() error {
	return c.doGet("login")
}

func (c *Client) Logout() error {
	return c.doGet("logout")
}

func (c *Client) GetProject() error {
	return c.doGet("project")
}

func (c *Client) UpdateProject() error {
	return c.doPost("project")
}
