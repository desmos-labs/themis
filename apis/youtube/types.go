package youtube

import "time"

// Channel contains the details of a Youtube user
type Channel struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
}

// NewChannel allows to build a new Channel instance
func NewChannel(id, title, description string, publishedAt time.Time) Channel {
	return Channel{
		ID:          id,
		Title:       title,
		Description: description,
		PublishedAt: publishedAt,
	}
}
