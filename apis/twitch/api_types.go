package twitch

// tokenResponseJSON contains the data returned from the Twitch APIs when getting a new OAuth token
type tokenResponseJSON struct {
	AccessToken string `json:"access_token"`
}

// usersResponseJSON contains the data returned from the Twitch APIs when reading the details of multiple users
type usersResponseJSON struct {
	Data []userResponseJSON `json:"data"`
}

// userResponseJSON contains the data returned from the Twitch APIs about a single user
type userResponseJSON struct {
	ID          string `json:"id"`
	Username    string `json:"display_name"`
	Description string `json:"description"`
}
