package main

import "fmt"

func main() {
	const value = 42
	i := int(value)
	f := float64(value)
	fmt.Println("i = ", i)
	fmt.Println("f = ", f)
}
