package template

import (
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
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
			"customKey": 0,
		}).
		Apply(`My {{ .env.CAT }} cat fought {{ .customKey }} {{ .env.HORSE }} horses.`)
	assert.NoError(t, err)
	assert.Equal(t, "My SPOOKY cat fought 0 EVIL horses.", tmpl)
}

func TestInvalidTemplate(t *testing.T) {
	_, err := New(context.New(config.Project{})).Apply(`{{ .timestamp`)
	assert.Error(t, err)
}

func TestEmptyTemplate(t *testing.T) {
	tmpl, err := New(context.New(config.Project{})).Apply("")
	assert.NoError(t, err)
	assert.Empty(t, tmpl)
}
