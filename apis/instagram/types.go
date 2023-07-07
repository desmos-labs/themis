package instagram

// User contains all the data of a user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// NewUser builds a new User instance
func NewUser(id, username, biography string) *User {
	return &User{
		ID:       id,
		Username: username,
		Bio:      biography,
	}
}
