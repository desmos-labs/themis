package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var fields = []string{"items.id", "items.snippet.title", "items.snippet.description", "items.snippet.publishedAt"}

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
	return ioutil.ReadAll(resp.Body)
}

// GetChannel returns the details of the Channel having the given id
func (api *API) GetChannel(id string) (Channel, error) {
	// Build the request endpoint. Example:
	// https://www.googleapis.com/youtube/v3/channels?id=<channel-id>&part=snippet&fields=items.id,items.snippet.title,items.snippet.description,items.snippet.publishedAt
	endpoint := fmt.Sprintf("%s/channels?id=%s&part=snippet&fields=%s", api.endpoint, id, strings.Join(fields, ","))

	// Create the request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Channel{}, err
	}

	// Run the request
	bz, err := api.runRequest(req)
	if err != nil {
		return Channel{}, err
	}

	var response channelResponseJSON
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return Channel{}, err
	}
	return convertChannelResponse(response)
}

// convertChannelResponse converts the given channelResponseJSON into a Channel object
func convertChannelResponse(response channelResponseJSON) (Channel, error) {
	if len(response.Items) != 1 {
		return Channel{}, fmt.Errorf("invalid number of channels returned: %d", len(response.Items))
	}
	channel := response.Items[0]
	return NewChannel(
		channel.ID,
		channel.Snippet.Title,
		channel.Snippet.Description,
		channel.Snippet.PublishedAt,
	), nil
}
