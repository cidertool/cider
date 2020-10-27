package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cidertool/cider/pkg/config"
)

// ErrConfigNotFound happens if a config file could not be found at any of the default locations.
var ErrConfigNotFound = errors.New("config file not found at any default path")

func loadConfig(path string, wd string) (config.Project, error) {
	if path != "" {
		return config.Load(path)
	}

	for _, f := range [4]string{
		".cider.yml",
		".cider.yaml",
		"cider.yml",
		"cider.yaml",
	} {
		proj, err := config.Load(filepath.Join(wd, f))
		if err != nil && os.IsNotExist(err) {
			continue
		}

		return proj, err
	}

	return config.Project{}, ErrConfigNotFound
}
