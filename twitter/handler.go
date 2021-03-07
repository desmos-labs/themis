package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/desmos-labs/themis/utils"
	"io/ioutil"
	"strings"
)

// Handler allows to handle Twitter-related operations properly
type Handler struct {
	cacheFilePath string
	api           *Api
}

// NewHandler returns a new Handler instance
func NewHandler(cacheFilePath string, api *Api) *Handler {
	return &Handler{
		cacheFilePath: cacheFilePath,
		api:           api,
	}
}

// cacheData represents how the Twitter-related data is stored inside the local cache
type cacheData struct {
	Tweets map[string]*Tweet // Maps the tweet id to the tweet
	Users  map[string]*User  // Maps the username to their user objects
}

// newCacheData returns a new empty cacheData instance
func newCacheData() *cacheData {
	return &cacheData{
		Tweets: map[string]*Tweet{},
		Users:  map[string]*User{},
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

// cacheTweet caches the given tweet for future references
func (h *Handler) cacheTweet(tweet *Tweet) error {
	cache, err := h.readCache()
	if err != nil {
		return err
	}

	// Set the tweet
	cache.Tweets[tweet.ID] = tweet

	// Serialize the contents
	bz, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	// Write the file
	return ioutil.WriteFile(h.cacheFilePath, bz, 0644)
}

// getTweetFromCache returns the tweet with the given id from the cache, if existing
func (h *Handler) getTweetFromCache(id string) (*Tweet, error) {
	cache, err := h.readCache()
	if err != nil {
		return nil, err
	}

	tweet, ok := cache.Tweets[id]
	if !ok {
		return nil, nil
	}
	return tweet, nil
}

// GetTweet returns the tweet having the given id, either from the cache if present of by querying the APIs.
// If the tweet was not cached, after retrieving it from the APIs it is later cached for future requests
func (h *Handler) GetTweet(id string) (*Tweet, error) {
	// Check the validity of the id
	if strings.ContainsRune(id, ',') {
		return nil, fmt.Errorf("invalid tweet id: %s", id)
	}

	// Try getting the tweet from the cache
	cached, err := h.getTweetFromCache(id)
	if err != nil {
		return nil, err
	}

	// Return the cached tweet if existing
	if cached != nil {
		return cached, nil
	}

	// If not cached, get the tweet from the APIs
	tweet, err := h.api.GetTweet(id)
	if err != nil {
		return nil, err
	}

	// Store into the cache
	err = h.cacheTweet(tweet)
	if err != nil {
		return nil, err
	}

	// Return the retrieved tweet
	return tweet, nil
}

// cacheUser caches the given user for future references
func (h *Handler) cacheUser(user *User) error {
	cache, err := h.readCache()
	if err != nil {
		return err
	}

	// Set the tweet
	cache.Users[user.Username] = user

	// Serialize the contents
	bz, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	// Write the file
	return ioutil.WriteFile(h.cacheFilePath, bz, 0644)
}

// getUserFromCache returns the User object associated with the user having the given username, if existing
func (h *Handler) getUserFromCache(username string) (*User, error) {
	cache, err := h.readCache()
	if err != nil {
		return nil, err
	}

	bio, ok := cache.Users[username]
	if !ok {
		return nil, nil
	}
	return bio, nil
}

// GetUser returns the bio of the user having the given username, either from the cache if present of
// by querying the APIs.
// If the user was not cached, after retrieving it from the APIs it is later cached for future requests
func (h *Handler) GetUser(username string) (*User, error) {
	// Check the validity of the username
	if strings.ContainsRune(username, ',') {
		return nil, fmt.Errorf("invalid username: %s", username)
	}

	// Check if the user is cached
	cached, err := h.getUserFromCache(username)
	if err != nil {
		return nil, err
	}

	// If cached, return
	if cached != nil {
		return cached, nil
	}

	// If not cached, get the user from the APIs
	user, err := h.api.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Cache the user
	err = h.cacheUser(user)
	if err != nil {
		return nil, err
	}

	// Return the retrieved bio
	return user, nil
}
