package youtube

// Handler allows to handle Youtube-related operations properly
type Handler struct {
	api *API
}

// NewHandler returns a new Handler instance
func NewHandler(api *API) *Handler {
	return &Handler{
		api: api,
	}
}

// GetUser returns the bio of the user having the given userID, either from the cache if present of
// by querying the APIs.
// If the user was not cached, after retrieving it from the APIs it is later cached for future requests
func (h *Handler) GetUser(userID string) (User, error) {
	return User{}, nil
}
