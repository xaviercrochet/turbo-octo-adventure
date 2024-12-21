package feed_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xaviercrochet/turbo-octo-adventure/pkg/net"
)

type FeedClient struct {
	hostname   string
	port       string
	httpClient *http.Client
}

func NewFeedClient(hostname, port string) *FeedClient {
	return &FeedClient{
		hostname:   hostname,
		port:       port,
		httpClient: &http.Client{},
	}
}

func (c *FeedClient) buildURL(path string) string {
	return fmt.Sprintf("http://%s:%s/api/%s", c.hostname, c.port, path)
}

/*
Call /api/healthz

return true if response is 200, false otherwise
*/
func (c *FeedClient) CheckHealth() (bool, error) {
	url := c.buildURL("healthz")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed creating http request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to query feed api: %w", err)
	}
	defer resp.Body.Close()

	if err := net.HttpStatusCodeToErr(resp); err != nil {
		return false, err
	}

	return true, nil
}

/*
Call POST "/select/feed"

params:
  - username: the username of the feed
  - accessToken: the access token

return the errors defined under pkg.net.errors based on the http status code of the response
*/
func (c *FeedClient) SelectFeed(selectedFeed, accessToken string) error {
	url := c.buildURL("select_feed")

	// build the payload for the post request
	payload := map[string]interface{}{
		"name": selectedFeed,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed serializing payload: %v\n", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed creating http request: %v\n", err)
	}

	// build the authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to query feed api: %v", err)
	}
	defer resp.Body.Close()

	// fail based on error code if not 200
	return net.HttpStatusCodeToErr(resp)
}

/*

Call /api/feed

params:
  - accessToken: the access token

  if successful, returns a list of songs

  return the errors defined under feed_api.errors based on the http status code of the response otherwise
*/

func (c *FeedClient) GetFeed(accessToken string) (*FeedResponse, error) {
	url := c.buildURL("feed")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating http request: %v", err)
	}

	// build the authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query feed api: %w", err)
	}
	defer resp.Body.Close()

	if err := net.HttpStatusCodeToErr(resp); err != nil {
		return nil, fmt.Errorf("failed to reach API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	var feedResponse FeedResponse
	err = json.Unmarshal(body, &feedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize json: %w", err)
	}

	// fail based on error code if not 200
	return &feedResponse, nil
}

type FeedResponse struct {
	WriteAccess bool  `json:"write_access"`
	Feed        *Feed `json:"feed"`
}

type Feed struct {
	Username string  `json:"username"`
	Songs    []*Song `json:"songs"`
}

type Song struct {
	Title      string    `json:"title"`
	ListenedAt time.Time `json:"listened_at"`
}
