package upperWaitGroup

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func doWork(uwg *UpperWaitGroup, count int, sleep time.Duration, done []bool) {
	for i := 0; i < count; i++ {
		if !uwg.Add() {
			break
		}
		go func(i int) {
			defer uwg.Done()

			done[i] = true
			time.Sleep(sleep)
		}(i)
	}
}

func TestUWG(t *testing.T) {
	uwg := NewUpperWaitGroup(5)

	count := 20
	done := make([]bool, count)

	doWork(uwg, count, time.Duration(0), done)

	uwg.Wait()

	for i := 0; i < count; i++ {
		if !done[i] {
			t.Errorf("done[%d] is not true", i)
		}
	}
}

func TestUWGWithCancel(t *testing.T) {
	uwg := NewUpperWaitGroup(5)

	count := 20
	done := make([]bool, count)

	go func() {
		time.Sleep(2 * time.Second)
		uwg.Cancel()
	}()

	doWork(uwg, count, time.Second, done)

	uwg.Wait()

	if !done[0] {
		t.Errorf("done[0] is not true")
	}

	if done[count-1] {
		t.Errorf("done[%d] is not false", count-1)
	}

}

// checking how many goroutines are running in parallel
// then increasing the parallelisation from 5 to 10
func TestUWGVariableUpperLimit(t *testing.T) {
	uwg := NewUpperWaitGroup(5)

	count := 120
	done := make([]bool, count)

	var currentParallel atomic.Int32
	var iterator atomic.Int32

	iterator.Store(-1)
	go func() {
		for iterator.Load() < int32(count) {
			cp := currentParallel.Load()
			i := iterator.Load()
			if i < 60 && cp > 5 {
				msg := fmt.Sprintf("[%d] Parallel running goroutines should be 4, but is: %d", i, cp)
				panic(msg)
			}
			if i > 80 && cp > 10 {
				msg := fmt.Sprintf("[%d] Parallel running goroutines should be at max 10, but is: %d", i, cp)
				panic(msg)
			}
			if i > 80 && i < 90 && cp < 9 {
				// after (count - upper_limit) there is not enough work for all workers
				msg := fmt.Sprintf("[%d] Parallel running goroutines should be at least 9, but is: %d", i, cp)
				panic(msg)
			}
		}
	}()
	for iterator.Add(1) < int32(count) {
		if !uwg.Add() {
			break
		}
		currentParallel.Add(1)

		if iterator.Load() == 60 {
			uwg.SetUpper(10)
		}
		go func(i int) {
			defer func() {
				uwg.Done()
				currentParallel.Add(-1)
			}()
			done[i] = true
			time.Sleep(time.Second / 10)
		}(int(iterator.Load()))
	}

	uwg.Wait()

	if !done[0] {
		t.Errorf("done[0] is not true")
	}

	if !done[count-1] {
		t.Errorf("done[count-1] is not true")
	}
}