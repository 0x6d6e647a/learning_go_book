package main

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	const numEntries int = 10_000_000
	// var entries []Person;
	entries := make([]Person, 0, numEntries)

	for i := 0; i < numEntries; i++ {
		entries = append(entries, Person{"Anthony", "Mendez", 36})
	}
}
