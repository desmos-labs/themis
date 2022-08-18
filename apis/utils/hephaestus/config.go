package hephaestus

// Config contains the configuration data for the Hephaestus integration
type Config struct {
	// StoreFolderPath represents the path to the folder inside which the data will be stored
	StoreFolderPath string `toml:"store_folder_path"`

	// HephaestusPubKeyPath represents the path to the file inside which is written the Hephaestus public key
	HephaestusPubKeyPath string `toml:"bot_pub_key_path"`
}
