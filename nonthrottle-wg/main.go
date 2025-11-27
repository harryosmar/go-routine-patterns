package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 30; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(1 * time.Second)
		}(i)
	}

	wg.Wait()
	log.Println("DONE")
}
