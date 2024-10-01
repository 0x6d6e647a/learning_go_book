package main

import "fmt"

func prefixer(prefix string) func(suffix string) string {
	return func(suffix string) string {
		return prefix + " " + suffix
	}
}

func main() {
	helloPrefix := prefixer("Hello")
	fmt.Println(helloPrefix("Bob"))
	fmt.Println(helloPrefix("Maria"))
}
