package twitter

import "time"

// userCacheData contains the details of a user cache. Each cache is supposed to last 1 hour.
type userCacheData struct {
	CacheTime time.Time
	User      *User
}

// newUserCacheData allows to build a new userCacheData for the given user
func newUserCacheData(user *User) *userCacheData {
	return &userCacheData{
		CacheTime: time.Now(),
		User:      user,
	}
}

// Expired tells whether this userCacheData contains expired data or not
func (c *userCacheData) Expired() bool {
	return c.CacheTime.Add(time.Hour).Before(time.Now())
}

// -----------------------------------------------------------------------------------------------------------------

// cacheData represents how the Twitter-related data is stored inside the local cache
type cacheData struct {
	// Maps the tweet id to the tweet
	Tweets map[string]*Tweet

	// Maps the username to their user objects
	Users map[string]*userCacheData
}

// newCacheData returns a new empty cacheData instance
func newCacheData() *cacheData {
	return &cacheData{
		Tweets: map[string]*Tweet{},
		Users:  map[string]*userCacheData{},
	}
}
