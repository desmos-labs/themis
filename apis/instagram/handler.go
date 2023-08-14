package instagram

import (
	"encoding/json"
	"os"

	"github.com/desmos-labs/themis/apis/utils"
)

// Handler allows to handle Instagram related requests
type Handler struct {
	api           *API
	cacheFilePath string
}

// NewHandler allows to build a new Handler instance
func NewHandler(cacheFilePath string, api *API) *Handler {
	return &Handler{
		cacheFilePath: cacheFilePath,
		api:           api,
	}
}

// cacheData represents how the Instagram data are stored inside the local cache
type cacheData struct {
	Medias map[string]*UserMedia // Maps the username to their media objects
}

// newCacheData returns a new empty cacheData instance
func newCacheData() *cacheData {
	return &cacheData{
		Medias: map[string]*UserMedia{},
	}
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

// cacheUser caches the given user for future references
func (h *Handler) cacheUser(user *UserMedia) error {
	cache, err := h.readCache()
	if err != nil {
		return err
	}

	// Set the media
	cache.Medias[user.Username] = user

	// Serialize the contents
	bz, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	// Write the file
	return os.WriteFile(h.cacheFilePath, bz, 0600)
}

// getUserMediaFromCache returns the Media object associated with the user having the given username, if existing
func (h *Handler) getUserMediaFromCache(username string) (*UserMedia, error) {
	cache, err := h.readCache()
	if err != nil {
		return nil, err
	}

	user, ok := cache.Medias[username]
	if !ok {
		return nil, nil
	}
	return user, nil
}

// GetUserMedia returns the media of the user having the given username from cache.
func (h *Handler) GetUserMedia(username string) (*UserMedia, error) {
	// Try getting the cached user media
	user, err := h.getUserMediaFromCache(username)
	if err != nil {
		return nil, err
	}

	// Return the user
	return user, nil
}

// RequestUserMedia requests the instagram user latest media from Instagram Graph API then store it inside cache.
func (h *Handler) RequestUserMedia(accessToken string) error {
	user, err := h.api.GetUserMedia(accessToken)
	if err != nil {
		return err
	}

	return h.cacheUser(user)
}
