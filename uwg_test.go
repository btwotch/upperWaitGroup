package upperWaitGroup

import (
	"testing"
	"time"
)

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

func doWork(uwg *upperWaitGroup, count int, sleep time.Duration, done []bool) {
	for i := 0; i < count; i++ {
		if !uwg.Add() {
			break
		}
		go func(i int) {
			defer func() {
				uwg.Done()
			}()

			done[i] = true
			time.Sleep(sleep)
		}(i)
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