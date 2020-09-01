package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterruptOK(t *testing.T) {
	assert.NoError(t, NewInterrupt().Run(context.Background(), func() error {
		return nil
	}))
}

func TestInterruptErrors(t *testing.T) {
	var err = errors.New("some error")
	assert.EqualError(t, NewInterrupt().Run(context.Background(), func() error {
		return err
	}), err.Error())
}

// func TestInterruptTimeout(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
// 	defer cancel()
// 	assert.EqualError(t, NewInterrupt().Run(ctx, func() error {
// 		t.Log("slow task...")
// 		time.Sleep(time.Minute)
// 		return nil
// 	}), "context deadline exceeded")
// }

// func TestInterruptSignals(t *testing.T) {
// 	for _, signal := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
// 		signal := signal
// 		t.Run(signal.String(), func(tt *testing.T) {
// 			tt.Parallel()
// 			var h = NewInterrupt()
// 			var errs = make(chan error, 1)
// 			go func() {
// 				errs <- h.Run(context.Background(), func() error {
// 					tt.Log("slow task...")
// 					time.Sleep(time.Minute)
// 					return nil
// 				})
// 			}()
// 			h.signals <- signal
// 			assert.EqualError(tt, <-errs, fmt.Sprintf("received: %s", signal))
// 		})
// 	}
// }
