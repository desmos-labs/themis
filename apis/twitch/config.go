package twitch

// Config contains the configuration details to properly use the Twitch APIs
type Config struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}
