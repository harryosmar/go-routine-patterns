package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	ch1, ch2 := broadcastChannel(makeChannel(data))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range ch1 {
			func(i int) {
				time.Sleep(1 * time.Second)
				log.Println("process A", i)
			}(i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range ch2 {
			func(i int) {
				time.Sleep(1 * time.Second)
				log.Println("process B", i)
			}(i)
		}
	}()

	wg.Wait()
	log.Println("Done")
}

func broadcastChannel(ch <-chan int) (<-chan int, <-chan int) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		for i := range ch {
			ch1 <- i
			ch2 <- i
		}
		close(ch1)
		close(ch2)
	}()
	return ch1, ch2
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
