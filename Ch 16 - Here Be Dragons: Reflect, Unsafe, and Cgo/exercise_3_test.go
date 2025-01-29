package ch16

import (
	"errors"
	"strings"
	"testing"
)

func TestParseOperation(t *testing.T) {
	testData := []struct {
		name        string
		ch          rune
		expectedOp  Operation
		expectedErr error
	}{
		{"add", '+', Addition, nil},
		{"sub", '-', Subtraction, nil},
		{"mul", '*', Multiplication, nil},
		{"div", '/', Division, nil},
		{"bad", '?', InvalidOp, ErrInvalidOpRune{'?'}},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			op, err := ParseOperation(d.ch)
			if op != d.expectedOp {
				t.Errorf("Expected operation '%v', got '%v'", d.expectedOp, op)
			}
			if err != d.expectedErr {
				t.Errorf("Expected error '%v', got '%v'", d.expectedErr, err)
			}

			if err != nil && !strings.HasPrefix(err.Error(), "invalid operation rune") {
				t.Error("Badly formatted error string")
			}
		})
	}
}

func TestOperationChar(t *testing.T) {
	testData := []struct {
		name        string
		op          Operation
		expectedCh  byte
		expectedErr error
	}{
		{"add", Addition, '+', nil},
		{"sub", Subtraction, '-', nil},
		{"mul", Multiplication, '*', nil},
		{"div", Division, '/', nil},
		{"bad", InvalidOp, 0, ErrInvalidOp},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			ch, err := d.op.Char()
			if ch != d.expectedCh {
				t.Errorf("Expected byte '%x', got '%x'", d.expectedCh, ch)
			}
			if !errors.Is(err, d.expectedErr) {
				t.Errorf("Expected error '%v', got '%v'", d.expectedErr, err)
			}
		})
	}
}

func TestOperationRun(t *testing.T) {
	testData := []struct {
		name        string
		op          Operation
		a           int
		b           int
		expectedRes int
		expectedErr error
	}{
		{"add", Addition, 1, 1, 2, nil},
		{"sub", Subtraction, 1, 1, 0, nil},
		{"mul", Multiplication, 1, 1, 1, nil},
		{"div", Division, 1, 1, 1, nil},
		{"bad", InvalidOp, 0, 0, 0, ErrInvalidOp},
		{"dbz", Division, 1, 0, 0, nil},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			res, err := d.op.Run(d.a, d.b)
			if res != d.expectedRes {
				t.Errorf("Expected result '%d', got '%d'", d.expectedRes, res)
			}
			if err != d.expectedErr {
				t.Errorf("Expected error '%v', got '%v'", d.expectedErr, err)
			}
		})
	}
}
