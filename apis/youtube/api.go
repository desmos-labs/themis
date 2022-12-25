package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	fields = strings.Join([]string{
		"items.id",
		"items.snippet.title",
		"items.snippet.description",
		"items.snippet.publishedAt",
	}, ",")
)

// API allows to query data from the Twitter APIs
type API struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

// NewAPI builds a new API instance
func NewAPI(apiKey string) *API {
	return &API{
		endpoint: "https://www.googleapis.com/youtube/v3",
		apiKey:   apiKey,
		client:   &http.Client{},
	}
}

// runRequest runs the given request, and returns the response body.
// If the request contains error, or the server answers with a status code different than 200,
// the method will return an error
func (api *API) runRequest(req *http.Request) ([]byte, error) {
	// Add the headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-goog-api-key", api.apiKey)

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
	return io.ReadAll(resp.Body)
}

func (api *API) getChannel(searchQuery string) (*Channel, error) {
	endpoint := fmt.Sprintf("%s/channels?%s&part=snippet&fields=%s", api.endpoint, searchQuery, fields)

	// Create the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Run the request
	bz, err := api.runRequest(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	var response channelResponseJSON
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	return convertChannelResponse(response)
}

// GetChannel returns the details of the Channel having the given id or username
func (api *API) GetChannel(search string) (*Channel, error) {
	// Try searching with the ID
	channel, err := api.getChannel(fmt.Sprintf("id=%s", search))
	if err != nil {
		return nil, err
	}

	// If not found, try searching with the username
	if channel == nil {
		channel, err = api.getChannel(fmt.Sprintf("forUsername=%s", search))
	}

	return channel, err
}

// convertChannelResponse converts the given channelResponseJSON into a Channel object
func convertChannelResponse(response channelResponseJSON) (*Channel, error) {
	if len(response.Items) == 0 {
		return nil, nil
	}

	if len(response.Items) > 1 {
		return nil, fmt.Errorf("invalid number of channels returned: %d", len(response.Items))
	}

	return NewChannel(
		response.Items[0].ID,
		response.Items[0].Snippet.Title,
		response.Items[0].Snippet.Description,
		response.Items[0].Snippet.PublishedAt,
	), nil
}
