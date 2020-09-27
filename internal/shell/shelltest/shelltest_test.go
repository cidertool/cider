package shelltest_test

import (
	"testing"

	"github.com/cidertool/cider/internal/shell/shelltest"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	ctx := context.New(config.Project{})
	ctx.CurrentDirectory = "TEST"
	sh := shelltest.Shell{
		Context: ctx,
		Commands: []shelltest.Command{
			{
				ReturnCode: 0,
				Stdout:     "TEST",
				Stderr:     "TEST",
			},
			{
				ReturnCode: 128,
				Stdout:     "TEST",
				Stderr:     "TEST",
			},
		},
	}

	dir := sh.CurrentDirectory()
	assert.Equal(t, dir, ctx.CurrentDirectory)

	exists := sh.Exists("echo")
	assert.True(t, exists)
	sh.SupportedPrograms = map[string]bool{
		"echo": false,
	}
	exists = sh.Exists("echo")
	assert.False(t, exists)

	cmd := sh.NewCommand("echo", "true")
	assert.NotNil(t, cmd)
	proc, err := sh.Exec(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, proc)
	proc, err = sh.Exec(cmd)
	assert.EqualError(t, err, "128")
	assert.NotNil(t, proc)
}
