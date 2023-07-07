package instagram

// User contains all the data of a user
type User struct {
	ID  string `json:"id"`
	Bio string `json:"biography"`
}

// NewUser builds a new User instance
func NewUser(id, username, biography string) *User {
	return &User{
		ID:  id,
		Bio: biography,
	}
}
