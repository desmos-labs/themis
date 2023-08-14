package instagram

// userResponse contains the data that are returned from the API used to get the info of a user
type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Media    media  `json:"media"`
}

type media struct {
	Data []mediaData `json:"data"`
}

type mediaData struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
}
