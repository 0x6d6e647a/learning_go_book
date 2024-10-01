package main

import (
	"fmt"
	"math/rand"
)

func main() {
	vals := make([]int, 100)

	for i := range vals {
		vals[i] = rand.Intn(100)
	}

	fmt.Println(vals)
}
