package main

import "fmt"

func main() {
	done := make(chan struct{})
	close(done)
	go func() {
		s, ok := <-done
		fmt.Println(s, ok)
	}()

	s, ok := <-done
	fmt.Println(s, ok)
	s, ok = <-done
	fmt.Println(s, ok)
}
