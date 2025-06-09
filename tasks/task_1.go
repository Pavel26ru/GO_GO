package task1

import (
	"fmt"
	"sync"
)

func TaskRun1() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine")
	}()

	wg.Wait()
}
