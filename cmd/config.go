package cmd

import (
	"fmt"
	"os"

	"github.com/aaronsky/applereleaser/pkg/config"
)

func loadConfig(path string) (config.Project, error) {
	if path != "" {
		return config.Load(path)
	}
	for _, f := range [4]string{
		".applereleaser.yml",
		".applereleaser.yaml",
		"applereleaser.yml",
		"applereleaser.yaml",
	} {
		proj, err := config.Load(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return proj, err
	}
	fmt.Println("Could not find a config file, using defaults.")
	return config.Project{}, nil
}
