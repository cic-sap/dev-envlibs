package util

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestFanOut(t *testing.T) {
	cc := make(chan int, 100)
	ci := int32(0)
	f := NewFanOut(5, func() {
		for c := range cc {
			fmt.Printf("%d\n", c)
			atomic.AddInt32(&ci, 1)
		}
	})
	for i := 0; i < 100; i++ {
		cc <- i
	}
	close(cc)
	f.Wait()
	if ci != 100 {
		t.Fatal(ci)
	}
	t.Logf("%d", ci)
}
