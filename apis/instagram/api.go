package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	fields = strings.Join([]string{
		"id",
		"username",
		"media.limit(1){caption}",
	}, ",")
)

// API allows to query data from the Instagram APIs
type API struct {
	endpoint string
	client   *http.Client
}

// NewAPI allows to build a new API instance
func NewAPI() *API {
	return &API{
		endpoint: "https://graph.instagram.com",
		client:   &http.Client{},
	}
}

// GetUserMedia returns the latest user media by access token provided by the user
func (api *API) GetUserMedia(accessToken string) (*UserMedia, error) {
	// Build the endpoint
	endpoint := fmt.Sprintf("%s/me?fields=%s&access_token=%s", api.endpoint, fields, accessToken)

	// Build the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Perform the request and check the response status code
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	// Parse the body
	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var response userResponse
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Media.Data) == 0 {
		return nil, fmt.Errorf("failed to get user latest media")
	}

	// Return the user media
	return NewUserMedia(
		response.Username,
		response.Media.Data[0].Caption,
	), nil
}
