package upperWaitGroup

import (
	"sync"
	"sync/atomic"
)

type upperWaitGroup struct {
	waitGroup sync.WaitGroup
	doCancel  atomic.Bool
	current   atomic.Int32
	waitMutex sync.Mutex
	upper     atomic.Int32
}

func NewUpperWaitGroup(max int) *upperWaitGroup {
	var uwg upperWaitGroup

	uwg.current.Store(0)
	uwg.upper.Store(int32(max))

	return &uwg
}

func (uwg *upperWaitGroup) Add() bool {
	if uwg.doCancel.Load() {
		return false
	}
	uwg.waitGroup.Add(1)
	for {
		if uwg.current.Add(1) <= uwg.upper.Load() {
			return true
		}
		uwg.current.Add(-1)
		uwg.waitMutex.Lock()
	}
}

func (uwg *upperWaitGroup) SetUpper(max int) {
	uwg.upper.Store(int32(max))
}

func (uwg *upperWaitGroup) Done() {
	uwg.waitGroup.Done()
	uwg.current.Add(-1)
	uwg.waitMutex.TryLock()
	uwg.waitMutex.Unlock()
}

func (uwg *upperWaitGroup) Wait() {
	uwg.waitGroup.Wait()
}

func (uwg *upperWaitGroup) Cancel() {
	uwg.doCancel.Store(true)
}
