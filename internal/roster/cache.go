package roster

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"shifts-go/cli/ui"
	"shifts-go/internal/helper"
)

var dir = filepath.Join(os.Getenv("HOME"), ".cache", "shifts-go")
var lifetimeReportsDir = filepath.Join(os.Getenv(("HOME")), ".cache", "shifts-go", "lifetime")

type CacheResponse struct {
	LoginToken string `json:"loginToken"`
	ExpiryDate int64  `json:"expiryDate"`
}

func Read(name string) ([]byte, bool, error) {
	return readFile(path(name))
}

func ReadLifetime(name string) ([]byte, bool, error) {
	return readFile(lifetimePath(name))
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

func WriteLifetimeReport(name string, report Report) error {
	name = helper.NormaliseName(name)
	if err := ensureLifetimeDir(); err != nil {
		return err
	}
	contents, err := json.Marshal(report)
	if err != nil {
		return err
	}

	if err := os.WriteFile(lifetimePath(name), []byte(contents), 0644); err != nil {
		return err
	}

	fmt.Printf("Report for %v saved!\n", ui.BoldLightCyan.Render(report.Employee.Name))
	return nil
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

func LifeTimeReportsDir() string {
	return lifetimeReportsDir
}

func ensureDir() error {
	return os.MkdirAll(dir, 0755)
}

func ensureLifetimeDir() error {
	return os.MkdirAll(lifetimeReportsDir, 0755)
}

func path(name string) string {
	return filepath.Join(dir, name)
}

func lifetimePath(name string) string {
	return filepath.Join(lifetimeReportsDir, name)
}

func readFile(path string) ([]byte, bool, error) {
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
