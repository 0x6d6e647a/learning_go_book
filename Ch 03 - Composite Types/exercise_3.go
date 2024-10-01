package main

import "fmt"

type Employee struct {
	firstName string
	lastName string
	id int
}

func main() {
	rob := Employee{"Rob", "Pike", 1}
	fmt.Println(rob)

	ken := Employee{
		firstName: "Ken",
		lastName: "Thompson",
		id: 2,
	}
	fmt.Println(ken)

	var me Employee
	me.firstName = "Anthony"
	me.lastName = "Mendez"
	me.id = 3
	fmt.Println(me)
}
