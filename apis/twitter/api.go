package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	userFields = "user.fields=created_at,description"
)

// API allows to query data from the Twitter APIs
type API struct {
	endpoint string
	bearer   string
	client   *http.Client
}

// NewAPI builds a new API instance
func NewAPI(bearer string) *API {
	return &API{
		endpoint: "https://api.twitter.com",
		bearer:   bearer,
		client:   &http.Client{},
	}
}

// runRequest runs the given request, and returns the response body.
// If the request contains error, or the server answers with a status code different than 200,
// the method will return an error
func (api *API) runRequest(req *http.Request) ([]byte, error) {
	// Add the headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", api.bearer))

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

// GetUser returns the user associated with the given username
func (api *API) GetUser(username string) (*User, error) {
	// Build the request endpoint. Example:
	// https://api.twitter.com/2/users/by?usernames=twitterdev&user.fields=description
	endpoint := fmt.Sprintf("%s/2/users/by?usernames=%s&%s", api.endpoint, username, userFields)

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

	// Parse the body
	var response userResponseJSON
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	return convertUserResponse(response)
}

// convertUserResponse converts the given userResponseJSON into a User object
func convertUserResponse(response userResponseJSON) (*User, error) {
	if len(response.Data) != 1 {
		return nil, fmt.Errorf("invalid number of users returned: %d", len(response.Data))
	}
	user := response.Data[0]

	return NewUser(
		user.ID,
		user.Name,
		user.Username,
		user.Bio,
	), nil
}

// GetTweet returns the details of the Tweet having the given id
func (api *API) GetTweet(id string) (*Tweet, error) {
	// Build the request endpoint. Example:
	// https://api.twitter.com/2/tweets?ids=1228393702244134912&tweet.fields=text
	endpoint := fmt.Sprintf(
		"%s/2/tweets?ids=%s&tweet.fields=text&expansions=author_id&%s",
		api.endpoint, id, userFields)

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

	var response tweetResponseJSON
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	return convertTweetResponse(response)
}

// convertTweetResponse converts the given tweetResponseJSON to a Tweet object
func convertTweetResponse(response tweetResponseJSON) (*Tweet, error) {
	if len(response.Data) != 1 {
		return nil, fmt.Errorf("invalid number of tweets returned: %d", len(response.Data))
	}
	tweet := response.Data[0]

	var user *userJSON
	for i, u := range response.Includes.Users {
		if u.ID == tweet.AuthorID {
			user = &response.Includes.Users[i]
		}
	}

	if user == nil {
		return nil, fmt.Errorf("tweet author not presend inside includes")
	}

	return NewTweet(
		tweet.ID,
		tweet.Text,
		tweet.CreatedAt,
		NewUser(
			user.ID,
			user.Name,
			user.Username,
			user.Bio,
		),
	), nil
}
