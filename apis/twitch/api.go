package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API allows to query data from the Twitch APIs
type API struct {
	endpoint     string
	clientID     string
	clientSecret string
	client       *http.Client
}

// NewAPI builds a new API instance
func NewAPI(clientID, clientSecret string) *API {
	return &API{
		endpoint:     "https://api.twitch.tv",
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       &http.Client{},
	}
}

func (api *API) runRequest(req *http.Request) ([]byte, error) {
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
	return ioutil.ReadAll(resp.Body)
}

// getAuthToken returns the token needed to authorize the following request
func (api *API) getAuthToken() (string, error) {
	endpoint := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials&scope=%s",
		api.clientID,
		api.clientSecret,
		"user_read",
	)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}

	body, err := api.runRequest(req)
	if err != nil {
		return "", err
	}

	var res tokenResponseJSON
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (api *API) GetUser(username string) (*User, error) {
	endpoint := fmt.Sprintf("%s/helix/users?login=%s", api.endpoint, username)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	authToken, err := api.getAuthToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Client-Id", api.clientID)

	body, err := api.runRequest(req)
	if err != nil {
		return nil, err
	}

	var res usersResponseJSON
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("invalid username: %s", username)
	}

	user := res.Data[0]
	return NewUser(user.ID, user.Username, user.Description), nil
}
