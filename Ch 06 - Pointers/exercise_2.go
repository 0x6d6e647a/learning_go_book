package main

import "fmt"

func UpdateSlice(slice []string, string string) {
	slice[len(slice)-1] = string
	fmt.Println("UpdateSlice = ", slice)
}

func GrowSlice(slice []string, string string) []string {
	newSlice := append(slice, string)
	fmt.Println("GrowSlice = ", newSlice)
	return newSlice
}

func main() {
	shapes := []string{"octahedron", "hexahedron", "cone"}
	fmt.Println("shapes = ", shapes)
	UpdateSlice(shapes, "tetrahedron")
	fmt.Println("shapes = ", shapes)
	shapes = GrowSlice(shapes, "sphere")
	fmt.Println("shapes = ", shapes)
}
