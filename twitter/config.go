package twitter

type Config struct {
	Bearer        string `toml:"bearer"`
	CacheFilePath string `toml:"cache_file"`
}
