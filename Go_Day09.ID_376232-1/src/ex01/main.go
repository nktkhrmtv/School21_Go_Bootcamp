package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func crawlWeb(ctx context.Context, urls <-chan string) <-chan string {
	results := make(chan string) 
	var wg sync.WaitGroup       
	semaphore := make(chan struct{}, 8) 

	go func() {
		defer close(results) 

		for url := range urls {
			select {
			case <-ctx.Done(): 
				return
			case semaphore <- struct{}{}: 
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					defer func() { <-semaphore }() 

					body, err := fetchURL(ctx, url)
					if err != nil {
						fmt.Printf("Ошибка при запросе %s: %v\n", url, err)
						return
					}

					select {
					case results <- body:
					case <-ctx.Done(): 
						return
					}
				}(url)
			}
		}
		wg.Wait() 
	}()
	
	return results
}

func fetchURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	return string(body), nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения. Останавливаюсь...")
		cancel() 
	}()

	urls := make(chan string)
	go func() {
		defer close(urls)
		defer fmt.Println("Все URL записаны\n")
		for _, url := range []string{
			"https://example.com",
			"https://golang.org",
			"https://github.com",
			"https://stackoverflow.com",
			"https://google.com",
			"https://amazon.com",
			"https://reddit.com",
			"https://wikipedia.org",
			"https://yandex.ru",
		} {
			select {
			case urls <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	results := crawlWeb(ctx, urls)
	for result := range results {
		fmt.Println("Получен результат:", result[:100]) 
		fmt.Println()
	}

	fmt.Println("Программа завершена.")
}