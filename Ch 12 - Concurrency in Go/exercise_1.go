package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	ch := make(chan int)

	wg.Add(2)
	// Get positive numbers.
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i += 1 {
			ch <- i
		}
	}()

	// Get negative numbers.
	go func() {
		defer wg.Done()
		for i := -100; i < 0; i += 1 {
			ch <- i
		}
	}()

	// Close channel when done.
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Count and sum the generated values.
	sum := 0
	count := 0

	for x := range ch {
		sum += x
		count += 1
	}

	// Check results.
	fmt.Printf("%-30s", "correct number of values: ")
	if count == 200 {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}

	fmt.Printf("%-30s", "correct sum of values: ")
	if sum == 0 {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}
}
