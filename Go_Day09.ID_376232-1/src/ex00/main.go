package main

import (
	"fmt"
	"sync"
	"time"
)

func sleepSort(numbers []int) <-chan int {
	out := make(chan int) 
	var wg sync.WaitGroup  

	for _, num := range numbers {
		wg.Add(1) 
		go func(n int) {
			defer wg.Done() 
			time.Sleep(time.Duration(n) * time.Second) 
			out <- n                                
		} (num)
	}

	go func() {
		wg.Wait()   
		close(out) 
	}()

	return out
}

func main() {
	numbers := []int{3, 1, 4, 1, 5, 8, 2, 6, 5, 3, 5}

	sortedChan := sleepSort(numbers)

	for num := range sortedChan {
		fmt.Println(num)
	}
}