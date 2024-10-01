package main

import (
	"fmt"
	"math"
)

func main() {
	var b byte = math.MaxUint8
	b += 1
	fmt.Println("b = ", b)

	var smallI int32 = math.MaxInt32
	smallI += 1
	fmt.Println("smallI = ", smallI)

	var bigI uint64 = math.MaxUint64
	bigI += 1
	fmt.Println("bigI = ", bigI)
}
