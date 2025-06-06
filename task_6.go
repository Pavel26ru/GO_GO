package main

import (
	"context"
	"fmt"
	"time"
)

func StartBatchProcessor(ctx context.Context, input <-chan int) {
	batch := make([]int, 0, 5)
	timer := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				fmt.Println("Processed batch (ctx cancelled):", batch)
			}
			return

		case val, ok := <-input:
			if !ok {
				if len(batch) > 0 {
					fmt.Println("Processed batch (channel closed):", batch)
				}
				return
			}

			batch = append(batch, val)

			if len(batch) == 5 {
				fmt.Println("Processed batch:", batch)
				batch = batch[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(2 * time.Second)
			}

		case <-timer.C:
			if len(batch) > 0 {
				fmt.Println("Processed batch (timeout):", batch)
				batch = batch[:0]
			}
			timer.Reset(2 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	input := make(chan int)

	go StartBatchProcessor(ctx, input)

	go func() {
		for i := 1; i <= 12; i++ {
			input <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(input)
	}()

	time.Sleep(8 * time.Second)
	fmt.Println("Main finished")
}
