package static

import (
	"strings"
	"testing"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestExampleConfig(t *testing.T) {
	_, err := config.LoadReader(strings.NewReader(ExampleConfig))
	assert.NoError(t, err)
}
