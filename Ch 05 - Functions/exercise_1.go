package main

import (
	"errors"
	"fmt"
)

type infixMathOperation func(a int, b int) (int, error)

func add(a int, b int) (int, error) { return a + b, nil }
func sub(a int, b int) (int, error) { return a - b, nil }
func mul(a int, b int) (int, error) { return a * b, nil }
func div(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
func mod(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a % b, nil
}

func symbolToInfixMathOperation(symbol rune) (infixMathOperation, error) {
	switch symbol {
	case '+':
		return add, nil
	case '-':
		return sub, nil
	case '*':
		return mul, nil
	case '/':
		return div, nil
	case '%':
		return mod, nil
	default:
		return nil, errors.New("unknown operation symbol")
	}
}

type expression struct {
	a  int
	op rune
	b  int
}

func main() {
	expressions := [...]expression{
		{2, '+', 3},
		{2, '-', 3},
		{2, '*', 3},
		{2, '/', 3},
		{2, '%', 3},
		{2, '/', 0},
		{3, '%', 0},
		{3, '_', 0},
	}

	for index, expression := range expressions {
		repr := fmt.Sprintf("%d: %d %c %d", index, expression.a, expression.op, expression.b)

		op, err := symbolToInfixMathOperation(expression.op)
		if err != nil {
			fmt.Printf("%s = ???\n", repr)
			continue
		}

		result, err := op(expression.a, expression.b)
		if err != nil {
			fmt.Printf("%s = err\n", repr)
			continue
		}

		fmt.Printf("%s = %d\n", repr, result)
	}
}
