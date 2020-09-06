package closer

import "testing"

type mockCloser struct {
	CloseWithError bool
}

func (c mockCloser) Close() error {
	return nil
}

func TestCloser(t *testing.T) {
	c := mockCloser{}
	Close(c)
}
