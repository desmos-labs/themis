package discord

type Config struct {
	// StoreFolderPath represents the path to the folder inside which the data will be stored
	StoreFolderPath string `toml:"store_folder_path"`
}
