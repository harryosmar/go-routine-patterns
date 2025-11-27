package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	strCh := make(chan string)
	intCh := make(chan int)
	doneCh := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Second)
		strCh <- "one"
		close(strCh)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		intCh <- 1
		close(intCh)
	}()

	go func() {
		wg.Wait()
		doneCh <- true
		close(doneCh)
	}()

	for {
		select {
		case str, ok := <-strCh:
			if ok {
				fmt.Println(str)
			}
		case i, ok := <-intCh:
			if ok {
				fmt.Println(i)
			}
		case <-doneCh:
			return
		default:
			fmt.Println("waiting")
			time.Sleep(1 * time.Second)
		}
	}
}
