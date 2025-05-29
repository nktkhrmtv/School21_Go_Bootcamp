package main

import (
	"fmt"
	"sync"
)

func multiplex(channels ...chan interface{}) chan interface{} {
    out := make(chan interface{})

    go func() {
        defer close(out)
        var wg sync.WaitGroup
        wg.Add(len(channels))
        
        for _, ch := range channels {
            go func(c chan interface{}) {
                defer wg.Done()
                for msg := range c {
                    out <- msg
                }
            }(ch)
        }
        
        wg.Wait() 
    }()

    return out
}

func main() {
	ch1 := make(chan interface{}) 
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	out := multiplex(ch1, ch2, ch3)

	go func() {
		defer close(ch1)
		for i := 1; i <= 5; i++ {
			ch1 <- fmt.Sprintf("ch1: %d", i)
		}
	}()

	go func() {
		defer close(ch2)
		for i := 6; i <= 10; i++ {
			ch2 <- fmt.Sprintf("ch2: %d", i)
		}
	}()

	go func() {
		defer close(ch3)
		for i := 10; i <= 15; i++ {
			ch3 <- fmt.Sprintf("ch3: %d", i)
		}
	}()

	for val := range out {
		fmt.Println(val)
	}
}