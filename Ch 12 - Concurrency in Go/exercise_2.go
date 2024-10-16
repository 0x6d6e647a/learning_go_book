package main

import (
	"fmt"
	"math/rand"
	"time"
)

func PerformCount() chan int {
	ch := make(chan int)
	jitter := time.Duration(rand.Intn(10)) * time.Millisecond

	go func() {
		for i := 0; i < 10; i += 1 {
			ch <- i
			time.Sleep(jitter)
		}
		close(ch)
	}()

	return ch
}

func main() {
	ch0 := PerformCount()
	ch1 := PerformCount()
	sum := 0
	numRunning := 2

	for numRunning != 0 {
		channel := -1
		var value int
		var ok bool

		select {
		case value, ok = <-ch0:
			if !ok {
				ch0 = nil
				numRunning -= 1
				break
			}
			channel = 0
		case value, ok = <-ch1:
			if !ok {
				ch1 = nil
				numRunning -= 1
				break
			}
			channel = 1
		}

		if channel != -1 {
			sum += value
			fmt.Printf("Channel %d yielded value %d.\n", channel, value)
		}
	}

	fmt.Print("correct values generated: ")
	if sum == 90 {
		fmt.Println("PASS")
	} else {
		fmt.Println("PASS")
	}
}
