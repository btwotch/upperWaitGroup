package upperWaitGroup

import (
	"sync"
)

type upperWaitGroup struct {
	upperWg chan struct{}
	wg      sync.WaitGroup
}

func NewUpperWaitGroup(max int) *upperWaitGroup {
	var uwg upperWaitGroup

	uwg.upperWg = make(chan struct{}, max)

	return &uwg
}

func (uwg *upperWaitGroup) Add() {
	uwg.wg.Add(1)
	uwg.upperWg <- struct{}{}
}

func (uwg *upperWaitGroup) Done() {
	uwg.wg.Done()
	<-uwg.upperWg
}

func (uwg *upperWaitGroup) Wait() {
	uwg.wg.Wait()
}
