package util

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type  DoubleCacheLoad func() (v interface{}, err error)

type DoubleCache struct {
	d DoubleCacheLoad
	delay time.Duration
	index int
	data [2]interface{}
	mutex     *sync.RWMutex
	name string
}

func NewDoubleCache(d DoubleCacheLoad, delay time.Duration, name string) *DoubleCache {
	dc := &DoubleCache{
		d:     d,
		delay: delay,
		index: 0,
		mutex: new(sync.RWMutex),
		name: name,
	}
	var t func()
	t = func() {
		defer func() {
			if err := recover(); err != nil {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				log.Println("[DoubleCache]func recover", dc.name, err, string(buf[:n]))
			}
			time.AfterFunc(delay, t)
		}()
		t1 := time.Now()
		data, err  := dc.d()
		if err == nil {
			dc.mutex.Lock()
			index := dc.index ^ 1
			dc.data[index] = data
			dc.index = index
			dc.mutex.Unlock()
			t2 := time.Now()
			useTime := t2.Sub(t1).Seconds()
			log.Printf("[DoubleCache]switch %d %s %f\n", index, dc.name, useTime)
		} else {
			log.Printf("[DoubleCache]err %s %s\n", dc.name, err.Error())
		}
	}
	t()
	return dc
}

func (dc *DoubleCache) Get() (v interface{}, err error) {
	dc.mutex.RLock()
	v = dc.data[dc.index]
	dc.mutex.RUnlock()
	if v == nil {
		err = fmt.Errorf("double cache data nil")
	}
	return
}


