package youtube

import "time"

type channelResponseJSON struct {
	Items []itemJSON `json:"items"`
}

type itemJSON struct {
	ID      string      `json:"id"`
	Snippet snippetJSON `json:"snippet"`
}

type snippetJSON struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
}
