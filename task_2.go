package main

import (
	"fmt"
	"sync"
)

func Gorutine(in int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello from goroutine", in)
}

func main() {
	wg := &sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Gorutine(i, wg)
	}

	wg.Wait()
}
