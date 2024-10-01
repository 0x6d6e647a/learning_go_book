package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func MakePerson(firstName string, lastName string, age int) Person {
	return Person{firstName, lastName, age}
}

func MakePersonPointer(firstName string, lastName string, age int) *Person {
	return &Person{firstName, lastName, age}
}

func main() {
	person := MakePerson("Rob", "Pike", 68)
	personPtr := MakePersonPointer("Ken", "Thompson", 81)
	fmt.Println(person)
	fmt.Println(*personPtr)
}
