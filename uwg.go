package upperWaitGroup

import (
	"sync"
	"sync/atomic"
)

type upperWaitGroup struct {
	upperWg  chan struct{}
	wg       sync.WaitGroup
	doCancel atomic.Bool
}

func NewUpperWaitGroup(max int) *upperWaitGroup {
	var uwg upperWaitGroup

	uwg.upperWg = make(chan struct{}, max)

	return &uwg
}

func (uwg *upperWaitGroup) Add() bool {
	if uwg.doCancel.Load() {
		return false
	}
	uwg.wg.Add(1)
	uwg.upperWg <- struct{}{}
	return true
}

func (uwg *upperWaitGroup) Done() {
	uwg.wg.Done()
	<-uwg.upperWg
}

func (uwg *upperWaitGroup) Wait() {
	uwg.wg.Wait()
}

func (uwg *upperWaitGroup) Cancel() {
	uwg.doCancel.Store(true)
}
