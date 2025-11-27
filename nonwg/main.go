package main

import (
	"fmt"
	"time"
)

func main() {
	output1 := make(chan int)
	output2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		output1 <- 1
		close(output1)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		output2 <- "two"
		close(output2)
	}()

	s := <-output2
	i := <-output1
	fmt.Println(i, s)
}
