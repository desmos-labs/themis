package youtube

import "time"

// channelCacheData contains the details of a user cache. Each cache is supposed to last 1 hour.
type channelCacheData struct {
	CacheTime time.Time
	Channel   Channel
}

// newChannelCacheData allows to build a new channelCacheData for the given channel
func newChannelCacheData(channel Channel) channelCacheData {
	return channelCacheData{
		CacheTime: time.Now(),
		Channel:   channel,
	}
}

// Expired tells whether this channelCacheData contains expired data or not
func (c channelCacheData) Expired() bool {
	return c.CacheTime.Add(time.Hour).Before(time.Now())
}

// -----------------------------------------------------------------------------------------------------------------

// cacheData represents how the Youtube-related data is stored inside the local cache
type cacheData struct {
	// Maps the username to their user objects
	Channels map[string]channelCacheData
}

// newCacheData returns a new empty cacheData instance
func newCacheData() *cacheData {
	return &cacheData{
		Channels: map[string]channelCacheData{},
	}
}
