package youtube

import (
	"net/http"
)

// API allows to query data from the Twitter APIs
type API struct {
	endpoint string
	client   *http.Client
}

// NewAPI builds a new API instance
func NewAPI() *API {
	return &API{
		endpoint: "https://api.youtube.com",
		client:   &http.Client{},
	}
}
