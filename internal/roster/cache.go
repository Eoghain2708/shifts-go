package roster

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var dir = filepath.Join(os.Getenv("HOME"), ".cache", "shifts-go")

type CacheResponse struct {
	LoginToken string `json:"loginToken"`
	ExpiryDate int64  `json:"expiryDate"`
}

func Read(name string) ([]byte, bool, error) {
	path := path(name)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if info.Size() == 0 {
		return nil, false, nil
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}

	return contents, true, nil
}

func Write(name string, cache CacheResponse) error {
	if err := ensureDir(); err != nil {
		return err
	}

	contents, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(path(name), []byte(contents), 0644)
}

func IsExist(name string) bool {
	_, err := os.Stat(path(name))
	return err == nil
}

func IsEmpty(name string) bool {
	info, err := os.Stat(path(name))
	if err != nil {
		return false
	}

	return info.Size() == 0
}

func ensureDir() error {
	return os.MkdirAll(dir, 0755)
}

func path(name string) string {
	return filepath.Join(dir, name)
}
