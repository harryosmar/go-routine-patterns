package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	iCh := percent(quad(random()))
	for i := range iCh {
		fmt.Printf("%.2f\n", i)
	}
}

func random() <-chan float64 {
	outputChan := make(chan float64)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			outputChan <- float64(i + 1)
		}
		close(outputChan)
	}()
	return outputChan
}

func quad(inputCh <-chan float64) <-chan float64 {
	outputChan := make(chan float64)
	go func() {
		for i := range inputCh {
			time.Sleep(1 * time.Second)
			outputChan <- i * i
		}
		close(outputChan)
	}()
	return outputChan
}

func percent(inputCh <-chan float64) <-chan float64 {
	outputChan := make(chan float64)
	go func() {
		for i := range inputCh {
			time.Sleep(1 * time.Second)
			outputChan <- i / 100
		}
		close(outputChan)
	}()
	return outputChan
}
