package main

import (
	"fmt"
	"math/rand"
)

func main() {
	vals := make([]int, 100)

	for i := range vals {
		n := &vals[i]
		*n = rand.Intn(100)

		switch {
		case *n%6 == 0:
			fmt.Println("Six!")
		case *n%2 == 0:
			fmt.Println("Two!")
		case *n%3 == 0:
			fmt.Println("Three!")
		default:
			fmt.Println("Never mind!")
		}
	}

}
