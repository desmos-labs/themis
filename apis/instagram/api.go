package instagram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API allows to query .data from the Instagram APIs
type API struct {
	endpoint string
	client   *http.Client
}

// NewAPI allows to build a new API instance
func NewAPI() *API {
	return &API{
		endpoint: "https://www.instagram.com",
		client:   &http.Client{},
	}
}

// GetUser returns the User associated to the given username
func (api *API) GetUser(username string) (*User, error) {
	// Build the endpoint
	endpoint := fmt.Sprintf("%s/%s/?__a=1", api.endpoint, username)

	// Build the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add the headers
	req.Header.Add("Content-Type", "application/json")

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
	bz, err := ioutil.ReadAll(resp.Body)
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
		response.GraphQL.User.Username,
		response.GraphQL.User.FullName,
		response.GraphQL.User.Biography,
	), nil
}
