// Package closer is a helper for closing io.Closers.
package closer

import (
	"io"
	"log"
)

type closeErrFunc func(error)

// nolint:gocheckglobals
var onCloseErr closeErrFunc = func(err error) { log.Fatal(err) }

// Close closes an io.Closer and handles the possible Close error.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		onCloseErr(err)
	}
}
