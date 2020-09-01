package template

import (
	"testing"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tmpl, err := New(context.New(config.Project{})).
		WithEnv(map[string]string{
			"DOG": "HAPPY",
			"CAT": "GRUMPY",
		}).
		WithShellEnv("HORSE=EVIL", "CAT=SPOOKY").
		WithFields(Fields{
			projectNameKey: "My Project",
			"customKey":    0,
		}).
		Apply(`{{ .projectName }}: My {{ .env.CAT }} cat fought {{ .customKey }} {{ .env.HORSE }} horses.`)
	assert.NoError(t, err)
	assert.Equal(t, "My Project: My SPOOKY cat fought 0 EVIL horses.", tmpl)
}

func TestInvalidTemplate(t *testing.T) {
	_, err := New(context.New(config.Project{})).
		WithFields(Fields{
			projectNameKey: "My Project",
		}).
		Apply(`{{ .projectName`)
	assert.Error(t, err)
}
