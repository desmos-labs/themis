package instagram

type userJSON struct {
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Biography string `json:"biography"`
}

// userResponse contains the data that are returned from the API used to get the info of a user
type userResponse struct {
	GraphQL struct {
		User userJSON `json:"user"`
	} `json:"graphql"`
}
