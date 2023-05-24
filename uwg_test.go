package upperWaitGroup

import "testing"

func TestUWG(t *testing.T) {
	uwg := NewUpperWaitGroup(5)

	count := 20
	done := make([]bool, count)

	for i := 0; i < count; i++ {
		uwg.Add()
		go func(i int) {
			defer func() {
				uwg.Done()
			}()

			done[i] = true
		}(i)
	}

	uwg.Wait()

	for i := 0; i < count; i++ {
		if !done[i] {
			t.Errorf("done[%d] is not true", i)
		}
	}
}
