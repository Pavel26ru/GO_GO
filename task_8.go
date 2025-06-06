package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- struct {
	url     string
	content string
}, client *http.Client, cancelFunc context.CancelFunc) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-jobs:
			if !ok {
				return
			}
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				cancelFunc()
				results <- struct {
					url, content string
				}{url, "error:" + err.Error()}
				return
			}

			resp, err := client.Do(req)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				cancelFunc()
				results <- struct {
					url, content string
				}{url, "error:" + err.Error()}
				return
			}
			defer resp.Body.Close()

			runes := []rune(string(bodyBytes))
			if len(runes) > 100 {
				runes = runes[:100]
			}

			results <- struct {
				url, content string
			}{url, string(runes)}
		}
	}
}

func FetchURLs(urls []string) map[string]string {
	const numJobs = 10
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	client := &http.Client{Timeout: 5 * time.Second}
	jobs := make(chan string)
	results := make(chan struct {
		url, content string
	}, len(urls))

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go worker(ctx, wg, jobs, results, client, cancel)
	}

	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	result := make(map[string]string)
	for res := range results {
		result[res.url] = res.content
	}

	return result
}

func main() {

	urls := []string{
		"https://example.org",
		"https://example.net",
		"https://example.edu",
		"https://example.tv",
		"https://example.cc",
	}

	result := FetchURLs(urls)

	fmt.Println(result)
}
