package task3

import (
	"fmt"
	"sync"
)

func TaskRun3() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go func() {
			defer wg.Done()
			ch <- i
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	sum := 0
	for val := range ch {
		sum += val
	}

	fmt.Println("Сумма:", sum)
}
