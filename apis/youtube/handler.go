package youtube

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/desmos-labs/themis/apis/utils"
)

// Handler allows to handle Youtube-related operations properly
type Handler struct {
	api           *API
	cacheFilePath string
}

// NewHandler returns a new Handler instance
func NewHandler(api *API, cacheFilePath string) *Handler {
	return &Handler{
		api:           api,
		cacheFilePath: cacheFilePath,
	}
}

// GetChannel returns the description of the user having the given user id
func (h *Handler) GetChannel(userID string) (Channel, error) {
	// Check the validity of the id
	if strings.ContainsRune(userID, ',') {
		return Channel{}, fmt.Errorf("invalid user id: %s", userID)
	}
	id := getChannelID(userID)
	// Try getting the channel from the cache
	cached, found, err := h.getChannelFromCache(id)
	if err != nil {
		return Channel{}, err
	}
	// Return the cached channel if existing
	if found {
		return cached, nil
	}
	// If not cached, get the channel from the APIs
	channel, err := h.api.GetChannel(id)
	if err != nil {
		return Channel{}, err
	}
	// Store into the cache
	err = h.cacheChannel(channel)
	if err != nil {
		return Channel{}, err
	}
	return channel, nil
}

// getChannelID returns the channel id from the given user id
// channel id is "UC" + <user id>
func getChannelID(userID string) string {
	return "UC" + userID
}

// readCache returns the current instance of the cache
func (h *Handler) readCache() (*cacheData, error) {
	bz, err := utils.ReadOrCreateFile(h.cacheFilePath)
	if err != nil {
		return nil, err
	}
	// Check if the file is empty
	if len(bz) == 0 {
		return newCacheData(), nil
	}
	// Deserialize the contents
	var data cacheData
	err = json.Unmarshal(bz, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// cacheChannel caches the given channel for future references
func (h *Handler) cacheChannel(channel Channel) error {
	cache, err := h.readCache()
	if err != nil {
		return err
	}
	// Set the channel
	cache.Channels[channel.ID] = newChannelCacheData(channel)
	return utils.WriteFile(h.cacheFilePath, cache)
}

// getChannelFromCache returns the channel with the given id from the cache, if existing
func (h *Handler) getChannelFromCache(id string) (Channel, bool, error) {
	cache, err := h.readCache()
	if err != nil {
		return Channel{}, false, err
	}

	data, ok := cache.Channels[id]
	if !ok {
		return Channel{}, false, nil
	}
	// If the data is expired, delete it
	if data.Expired() {
		delete(cache.Channels, id)
		return Channel{}, false, utils.WriteFile(h.cacheFilePath, cache)
	}
	return data.Channel, true, nil
}
