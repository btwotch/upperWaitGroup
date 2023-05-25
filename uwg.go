package upperWaitGroup

import (
	"sync"
	"sync/atomic"
)

type UpperWaitGroup struct {
	waitGroup sync.WaitGroup
	doCancel  atomic.Bool
	current   atomic.Int32
	waitMutex sync.Mutex
	upper     atomic.Int32
}

func NewUpperWaitGroup(max int) *UpperWaitGroup {
	var uwg UpperWaitGroup

	uwg.current.Store(0)
	uwg.upper.Store(int32(max))

	return &uwg
}

func (uwg *UpperWaitGroup) Add() bool {
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

func (uwg *UpperWaitGroup) SetUpper(max int) {
	uwg.upper.Store(int32(max))
}

func (uwg *UpperWaitGroup) Done() {
	uwg.waitGroup.Done()
	uwg.current.Add(-1)
	uwg.waitMutex.TryLock()
	uwg.waitMutex.Unlock()
}

func (uwg *UpperWaitGroup) Wait() {
	uwg.waitGroup.Wait()
}

func (uwg *UpperWaitGroup) Cancel() {
	uwg.doCancel.Store(true)
}
