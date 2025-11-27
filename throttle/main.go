package main

import (
	"log"
	"time"
)

func main() {
	sem := make(chan struct{}, 3)
	for i := 1; i <= 30; i++ {
		sem <- struct{}{}
		go func(i int) {
			defer func() { <-sem }()
			log.Println(i)
			time.Sleep(1 * time.Second)
		}(i)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	log.Println("DONE")
}
