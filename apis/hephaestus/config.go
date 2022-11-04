package hephaestus

// Config contains the configuration data for the Hephaestus integration
type Config struct {
	// PubKeyPath represents the path to the file inside which is written the Hephaestus public key
	PubKeyPath string `toml:"pubkey_path"`
}
