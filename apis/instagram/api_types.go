package instagram

// userResponse contains the data that are returned from the API used to get the info of a user
type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Biography string `json:"biography"`
}
