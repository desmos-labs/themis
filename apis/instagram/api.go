package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GetUser returns the User by access token provided by the user
func (api *API) GetUser(accessToken string) (*User, error) {
	// Build the endpoint
	endpoint := fmt.Sprintf("%s/me?fields=id,biography&accessToken=%s", api.endpoint, accessToken)

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

	// Return the user
	return NewUser(
		response.ID,
		response.Username,
		response.Biography,
	), nil
}
