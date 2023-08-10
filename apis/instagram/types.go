package instagram

// UserMedia contains media data
type UserMedia struct {
	Username string `json:"username"`
	Caption  string `json:"caption"`
}

// NewUserMedia builds a new UserMedia instance
func NewUserMedia(username, biography string) *UserMedia {
	return &UserMedia{
		Username: username,
		Caption:  biography,
	}
}

// AddUserMediaRequest represents the access token that allows themis to ask user media from instagram
type AddUserMediaRequest struct {
	AccessToken string `json:"accessToken"`
}
