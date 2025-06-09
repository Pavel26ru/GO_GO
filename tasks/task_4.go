package task4

import (
	"fmt"
	"sync"
)

var count int = 0

func countUp(mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	count++
	fmt.Println(count)
	mu.Unlock()
}
func TaskRun4() {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go countUp(mu, wg)
	}

	wg.Wait()
}
