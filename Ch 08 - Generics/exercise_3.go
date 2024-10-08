package main

import (
	"fmt"
	"strings"
)

type Node[T comparable] struct {
	Value T
	Next  *Node[T]
}

type List[T comparable] struct {
	Dummy Node[T]
	Length int
}

func (l *List[T]) Add(value T) {
	curr := &l.Dummy

	for curr.Next != nil {
		curr = curr.Next
	}

	curr.Next = &Node[T]{value, nil}
	l.Length += 1
}

func (l *List[T]) Insert(value T, index int) {
	curr := &l.Dummy

	for curr.Next != nil && index > 0 {
		curr = curr.Next
		index -= 1
	}

	curr.Next = &Node[T]{value, curr.Next}
	l.Length += 1
}

func (l *List[T]) Index(value T) int {
	index := 0
	curr := l.Dummy.Next

	for curr != nil {
		if curr.Value == value {
			return index
		}
		curr = curr.Next
		index += 1
	}

	return -1
}

func (l List[T]) String() string {
	var sb strings.Builder
	curr := l.Dummy.Next

	for curr != nil {
		repr := fmt.Sprint(curr.Value)
		sb.WriteString(repr)

		if curr.Next != nil {
			sb.WriteString(", ")
		}

		curr = curr.Next
	}

	return sb.String()
}

func main() {
	var list List[int]

	for index := range 4 {
		if index == 1 {
			continue
		}
		list.Add(index + 1)
	}

	list.Insert(0, 0)
	list.Insert(2, 2)

	fmt.Println(list)

	for index := range list.Length {
		if index != list.Index(index) {
			panic("insertion order to index sanity check failed")
		}
	}

	if list.Index(-42) != -1 {
		panic("index of non-existent value sanity check failed")
	}
}
