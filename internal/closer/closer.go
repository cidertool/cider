// Package closer contains a single utility for closing buffers in a defer
package closer

import (
	"io"
	"log"
)

// Close closes an open descriptor.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
