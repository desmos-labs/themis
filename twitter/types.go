package twitter

import "time"

// User contains the details of a Twitter user
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// NewUser allows to build a new User instance
func NewUser(id, name, username, bio string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Username: username,
		Bio:      bio,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// Tweet contains the details of a tweet
type Tweet struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	CreationTime time.Time `json:"creation_time"`
	Author       *User     `json:"author"`
}

// NewTweet allows to build a new Tweet instance
func NewTweet(id, text string, creationTime time.Time, author *User) *Tweet {
	return &Tweet{
		ID:           id,
		Text:         text,
		CreationTime: creationTime,
		Author:       author,
	}
}
