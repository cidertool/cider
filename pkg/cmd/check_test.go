package cmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCheckCmd(t *testing.T) {
	var cmd = newCheckCmd()

	var path = filepath.Join(t.TempDir(), "foo.yaml")

	var proj config.Project

	s, err := proj.String()
	assert.NoError(t, err)
	err = ioutil.WriteFile(path, []byte(s), 0600)
	assert.NoError(t, err)

	err = cmd.cmd.Execute()
	assert.Error(t, err)

	cmd.config = path

	err = cmd.cmd.Execute()
	assert.NoError(t, err)
}
