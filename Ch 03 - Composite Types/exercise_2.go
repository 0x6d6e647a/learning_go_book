package main

import "fmt"

func main() {
	message := "Hi 👩 and 👨"
	runes := []rune(message)
	fmt.Printf("%c\n", runes[3])
}
