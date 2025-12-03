package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := makeChannel([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for v := range ch {
				time.Sleep(1 * time.Second)
				fmt.Println("Group ", i, v)
			}
		}(i + 1)
	}
	wg.Wait()
	fmt.Println("done")
}

func makeChannel(data []int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, i := range data {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
