package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeSkip(t *testing.T) {
	skip := Skip("TEST")
	var err error = ErrSkip{reason: "TEST"}
	assert.Error(t, skip)
	assert.Error(t, err)
	assert.True(t, IsSkip(err))
	assert.Equal(t, err, skip)
	assert.Equal(t, err.Error(), skip.Error())
}
