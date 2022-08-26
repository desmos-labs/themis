package youtube

type Config struct {
	APIKEY        string `toml:"api_key"`
	CacheFilePath string `toml:"cache_file"`
}
