package parallel

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	var g = New(4)
	var lock sync.Mutex
	var counter int
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			lock.Lock()
			counter++
			lock.Unlock()
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	assert.Equal(t, counter, 10)
}

func TestGroupOrder(t *testing.T) {
	var num = 10
	var g = New(1)
	var output = []int{}
	for i := 0; i < num; i++ {
		i := i
		g.Go(func() error {
			output = append(output, i)
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, output)
}

func TestGroupOrderError(t *testing.T) {
	var g = New(1)
	var output = []int{}
	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			output = append(output, i)
			return fmt.Errorf("fake err")
		})
	}
	assert.EqualError(t, g.Wait(), "fake err")
	assert.Equal(t, []int{0}, output)
}
