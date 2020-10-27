package cmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Happy_CustomPath(t *testing.T) {
	var path = filepath.Join(t.TempDir(), "foo.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = ioutil.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)
	cfg, err := loadConfig(path, "")
	assert.NoError(t, err)
	assert.Empty(t, cfg)
}

func TestConfig_Happy_DefaultPath(t *testing.T) {
	var folder = t.TempDir()

	var path = filepath.Join(folder, "cider.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = ioutil.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)
	cfg, err := loadConfig("", folder)
	assert.NoError(t, err)
	assert.Empty(t, cfg)
}

func TestConfig_Err_DoesntExist(t *testing.T) {
	cfg, err := loadConfig("", t.TempDir())
	assert.Error(t, err)
	assert.Empty(t, cfg)
}
