package twitch

import (
	"fmt"
	"strings"
)

// Handler allows to handle Twitch-related operations properly
type Handler struct {
	api *API
}

// NewHandler returns a new Handler instance
func NewHandler(api *API) *Handler {
	return &Handler{
		api: api,
	}
}

// GetUser returns the bio of the user having the given username, either from the cache if present of
// by querying the APIs.
// If the user was not cached, after retrieving it from the APIs it is later cached for future requests
func (h *Handler) GetUser(username string) (*User, error) {
	// Check the validity of the username
	if strings.ContainsRune(username, ',') {
		return nil, fmt.Errorf("invalid username: %s", username)
	}

	// Get the user from the APIs
	user, err := h.api.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Return the retrieved bio
	return user, nil
}
