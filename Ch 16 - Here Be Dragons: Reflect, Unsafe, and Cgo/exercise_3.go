package ch16

/*
	extern int mini_calc(char*, int, int);
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type Operation uint8

const (
	InvalidOp Operation = iota
	Addition
	Subtraction
	Multiplication
	Division
)

type ErrInvalidOpRune struct {
	BadRune rune
}

func (e ErrInvalidOpRune) Error() string {
	return fmt.Sprintf("invalid operation rune '%c'", e.BadRune)
}

func ParseOperation(ch rune) (Operation, error) {
	switch ch {
	case '+':
		return Addition, nil
	case '-':
		return Subtraction, nil
	case '*':
		return Multiplication, nil
	case '/':
		return Division, nil
	}

	return InvalidOp, ErrInvalidOpRune{ch}
}

var ErrInvalidOp = errors.New("invalid operation")

func (op Operation) Char() (byte, error) {
	switch op {
	case Addition:
		return '+', nil
	case Subtraction:
		return '-', nil
	case Multiplication:
		return '*', nil
	case Division:
		return '/', nil
	default:
		return 0, ErrInvalidOp
	}
}

func (op Operation) Run(a int, b int) (int, error) {
	ch, err := op.Char()
	if err != nil {
		return 0, err
	}

	res := C.mini_calc(
		(*C.char)(unsafe.Pointer(&ch)),
		C.int(a),
		C.int(b),
	)

	return int(res), nil
}
