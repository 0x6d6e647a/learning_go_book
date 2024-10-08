package main

import "fmt"

type Printable interface {
	fmt.Stringer
	~int | ~float64
}

type PrintInt int

func (mi PrintInt) String() string {
	return fmt.Sprintf("%d", mi)
}

type PrintFloat float64

func (pf PrintFloat) String() string {
	return fmt.Sprintf("%0.2f", pf)
}

func PrintIt[T Printable](t T) {
	fmt.Println(t)
}

func main() {
	var i PrintInt = 10
	PrintIt(i)

	var f PrintFloat = 10.12
	PrintIt(f)
}
