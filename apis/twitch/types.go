package twitch

// User contains the details of a Twitter user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// NewUser allows to build a new User instance
func NewUser(id, username, bio string) *User {
	return &User{
		ID:       id,
		Username: username,
		Bio:      bio,
	}
}
