package main

import "fmt"

type Numberic interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Double[T Numberic](value T) T {
	return value * 2
}

func main() {
	fmt.Printf("%d\n", Double(2))
	fmt.Printf("%0.1f\n", Double(2.0))
}
