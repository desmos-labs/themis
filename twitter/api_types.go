package twitter

import "time"

// tweetJson contains the data of a tweet that is retrieved from the Twitter APIs
type tweetJson struct {
	AuthorId  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Text      string    `json:"text"`
}

// userJson contains the data of a user retrieved from the Twitter APIs
type userJson struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"description"`
}

// tweetResponseJson contains the data that is returned from the Twitter v2 /tweet REST API
type tweetResponseJson struct {
	Data     []tweetJson `json:"data"`
	Includes struct {
		Users []userJson `json:"users"`
	} `json:"includes"`
}
type userResponseJson struct {
	Data []userJson
}
