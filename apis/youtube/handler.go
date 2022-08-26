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

// GetChannel returns the description of the user having the given user id
func (h *Handler) GetChannel(userID string) (Channel, error) {
	// Check the validity of the id
	if strings.ContainsRune(userID, ',') {
		return Channel{}, fmt.Errorf("invalid user id: %s", userID)
	}
	return h.api.GetChannel(getChannelID(userID))
}

// getChannelID returns the channel id from the given user id
// channel id is "UC" + <user id>
func getChannelID(userID string) string {
	return "UC" + userID
}
