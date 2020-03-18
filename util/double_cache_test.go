package util

import (
	"log"
	"testing"
	"time"
)

func TestNewDoubleCache(t *testing.T) {
	i := 0
	dc := NewDoubleCache(func() (v interface{}, err error) {
		i++
		v = i
		return
	}, time.Millisecond * 100, "test")
	for ;; {
		time.Sleep(time.Millisecond * 50)
		index, err := dc.Get()
		if err != nil {
			log.Printf("test_data:%d, %s\n", index, err.Error())
		} else {
			log.Printf("test_data:%d\n", index)
		}
	}
}
