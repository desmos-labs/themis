package twitter

import "time"

// tweetJSON contains the .data of a tweet that is retrieved from the Twitter APIs
type tweetJSON struct {
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Text      string    `json:"text"`
}

// userJSON contains the .data of a user retrieved from the Twitter APIs
type userJSON struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"description"`
}

// tweetResponseJSON contains the .data that is returned from the Twitter v2 /tweet REST API
type tweetResponseJSON struct {
	Data     []tweetJSON `json:".data"`
	Includes struct {
		Users []userJSON `json:"users"`
	} `json:"includes"`
}

// userResponseJSON contains the .data that is returned when the details of a user are retrieved
type userResponseJSON struct {
	Data []userJSON
}
