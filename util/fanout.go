package util

import (
	"log"
	"sync"
)

type FanOut struct {
	n int
	wg *sync.WaitGroup
}

func NewFanOut(routineNum int, f func()) *FanOut {
	var wg sync.WaitGroup
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		go func() {
		   defer func() {
		   		if e := recover(); e != nil {
		   			log.Println("fan out recover", e)
				}
		   		wg.Done()
		   }()
		   f()
		}()
	}
	return &FanOut{
		n:  routineNum,
		wg: &wg,
	}
}

func(f *FanOut)Wait() {
	f.wg.Wait()
}
