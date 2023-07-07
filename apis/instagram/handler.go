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
	Users map[string]*User // Maps the username to their user objects
}

// newCacheData returns a new empty cacheData instance
func newCacheData() *cacheData {
	return &cacheData{
		Users: map[string]*User{},
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
func (h *Handler) cacheUser(user *User) error {
	cache, err := h.readCache()
	if err != nil {
		return err
	}

	// Set the tweet
	cache.Users[user.ID] = user

	// Serialize the contents
	bz, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	// Write the file
	return os.WriteFile(h.cacheFilePath, bz, 0600)
}

// getUserFromCache returns the User object associated with the user having the given user ID, if existing
func (h *Handler) getUserFromCache(userID string) (*User, error) {
	cache, err := h.readCache()
	if err != nil {
		return nil, err
	}

	user, ok := cache.Users[userID]
	if !ok {
		return nil, nil
	}
	return user, nil
}

// GetUser returns the bio of the user having the given username, either from the cache if present of
// by querying the APIs.
// If the user was not cached, after retrieving it from the APIs it is later cached for future requests
func (h *Handler) GetUser(username string) (*User, error) {
	// Try getting the cached user
	user, err := h.getUserFromCache(username)
	if err != nil {
		return nil, err
	}

	// Return the user
	return user, nil
}

// RequestUser requests the instagram user profile from Meta Graph API then store it inside cache.
func (h *Handler) RequestUser(userID string, accessToken string) error {
	user, err := h.api.GetUser(userID, accessToken)
	if err != nil {
		return err
	}

	return h.cacheUser(user)
}
