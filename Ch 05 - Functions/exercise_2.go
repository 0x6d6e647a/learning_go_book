package main

import (
	"fmt"
	"log"
	"os"
)

func fileLen(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no files specified")
	}

	for _, path := range os.Args[1:] {
		size, err := fileLen(path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d bytes :: %s\n", size, path)
	}
}
