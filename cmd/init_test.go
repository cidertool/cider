package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	var folder = t.TempDir()

	var cmd = newInitCmd().cmd

	var path = filepath.Join(folder, "foo.yaml")

	cmd.SetArgs([]string{"-f", path, "--skip-prompt"})
	assert.NoError(t, cmd.Execute())
	assert.FileExists(t, path)
}
