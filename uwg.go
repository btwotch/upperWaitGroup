package upperWaitGroup

import (
	"sync"
	"sync/atomic"
)

type UpperWaitGroup struct {
	waitGroup      sync.WaitGroup
	waitGroupMutex sync.Mutex
	doCancel       atomic.Bool
	current        atomic.Int32
	waitMutex      sync.Mutex
	upper          atomic.Int32
	doneMutex      sync.Mutex // prevent parallel calling of Done()
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
	uwg.waitGroupMutex.Lock()
	uwg.waitGroup.Add(1)
	uwg.waitGroupMutex.Unlock()
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

func (uwg *UpperWaitGroup) GetUpper() int {
	return int(uwg.upper.Load())
}

func (uwg *UpperWaitGroup) Done() {
	uwg.doneMutex.Lock()
	uwg.waitGroup.Done()
	uwg.current.Add(-1)

	uwg.waitMutex.TryLock()
	uwg.waitMutex.Unlock()
	uwg.doneMutex.Unlock()
}

func (uwg *UpperWaitGroup) Wait() {
	uwg.waitGroupMutex.Lock()
	uwg.waitGroup.Wait()
	uwg.waitGroupMutex.Unlock()
}

func (uwg *UpperWaitGroup) Cancel() {
	uwg.doCancel.Store(true)
}

func (uwg *UpperWaitGroup) Current() int {
	return int(uwg.current.Load())
}
