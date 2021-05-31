package instagram

// User contains all the data of a user
type User struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}

// NewUser builds a new User instance
func NewUser(username, fullName, biography string) *User {
	return &User{
		Username: username,
		FullName: fullName,
		Bio:      biography,
	}
}
