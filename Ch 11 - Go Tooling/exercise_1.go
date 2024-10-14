package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed english_rights.txt
var englishRights string

//go:embed deutsch_rechten.txt
var germanRights string

//go:embed 日本語_権.txt
var japaneseRights string

//go:embed español_derechos.txt
var spanishRights string

var langMap = map[string]*string{
	"english": &englishRights,
	"german": &germanRights,
	"japanese": &japaneseRights,
	"spanish": &spanishRights,
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("no language provided")
		os.Exit(1)
	}

	language := args[1]
	rightsText, ok := langMap[strings.ToLower(language)]
	if !ok {
		fmt.Println("unknown language: " + language)
		os.Exit(1)
	}

	fmt.Println(*rightsText)
}
