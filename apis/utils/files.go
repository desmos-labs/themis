package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	PermDir = 0777
	PerFile = 0666
)

// DoesFileExist tells whether or not the file located at the given path exists
func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDirIfNotExisting creates the parent directory of the given file path if not existing
func CreateDirIfNotExisting(filePath string) {
	// Create the parent directory if not existing
	cacheDir := filepath.Dir(filePath)
	if !DoesFileExist(cacheDir) {
		_ = os.Mkdir(cacheDir, PermDir)
	}
}

// GetFileOrCreate returns a pointer to the file located at the
// given path, creating it if it does not exist
func GetFileOrCreate(path string) (*os.File, error) {
	CreateDirIfNotExisting(path)

	// Open the file
	return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, PerFile)
}

// ReadOrCreateFile reads the file located at the given path, or creates it if not existing.
func ReadOrCreateFile(path string) ([]byte, error) {
	f, err := GetFileOrCreate(path)
	if err != nil {
		return nil, err
	}

	// Read the file contents
	return io.ReadAll(f)
}

// ReadFile reads the content of the file located at the given path.
// If the file does not exist, it returns an empty content.
func ReadFile(path string) ([]byte, error) {
	if !DoesFileExist(path) {
		return nil, nil
	}

	return os.ReadFile(os.ExpandEnv(path))
}

// WriteFile writes the given data inside the file located at the specified path
func WriteFile(path string, data interface{}) error {
	CreateDirIfNotExisting(path)

	// Serialize the contents
	bz, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	// Write the file
	return os.WriteFile(path, bz, PerFile)
}
