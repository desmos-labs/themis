package youtube

import (
	"fmt"
	"strings"
)

// Handler allows to handle Youtube-related operations properly
type Handler struct {
	api *API
}

// NewHandler returns a new Handler instance
func NewHandler(api *API) *Handler {
	return &Handler{
		api: api,
	}
}

// GetChannel returns the details of the user having the given username or id
func (h *Handler) GetChannel(search string) (*Channel, error) {
	// Validate the search value
	if strings.Contains(search, ",") {
		return nil, fmt.Errorf("invalid search value")
	}

	// Get the channel
	return h.api.GetChannel(search)
}
