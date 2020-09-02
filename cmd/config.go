package cmd

import (
	"os"
	"path/filepath"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/apex/log"
)

func loadConfig(path string, wd string) (config.Project, error) {
	if path != "" {
		return config.Load(path)
	}
	for _, f := range [4]string{
		".applereleaser.yml",
		".applereleaser.yaml",
		"applereleaser.yml",
		"applereleaser.yaml",
	} {
		proj, err := config.Load(filepath.Join(wd, f))
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return proj, err
	}
	log.Warn("could not find a config file, using defaults...")
	return config.Project{}, nil
}
