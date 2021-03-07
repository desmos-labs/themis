package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadOrCreateFile reads the file located at the given path, or creates it if not existing.
func ReadOrCreateFile(path string) ([]byte, error) {
	// Create the parent directory if not existing
	cacheDir := filepath.Dir(path)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		_ = os.Mkdir(cacheDir, 0644)
	}

	// Open the file
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the file contents
	return ioutil.ReadAll(f)
}
