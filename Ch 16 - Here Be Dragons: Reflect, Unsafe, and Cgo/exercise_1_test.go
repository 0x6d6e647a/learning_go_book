package ch16

import (
	"errors"
	"testing"
)

func TestValidateStringLength(t *testing.T) {
	// Non-struct rejection.
	err := ValidateStringLength(5)
	if !errors.Is(err, ErrNonStruct) {
		t.Error("non struct not rejected")
	}

	// Check skipping fields.
	t0 := struct {
		skipInt int
		skipStr string
	}{}
	err = ValidateStringLength(t0)
	if err != nil {
		t.Error("error skipping unchecked fields")
	}

	// Bad tag integer value.
	t1 := struct {
		test string `minStrLen:"zib"`
	}{}
	err = ValidateStringLength(t1)
	if err.Error() != "field 'test' minStrLen tag 'zib' is not an int" {
		t.Error("bad minStrLen value not detected")
	}

	// Test short string.
	t2 := struct {
		test string `minStrLen:"2"`
	}{"x"}
	err = ValidateStringLength(t2)
	if !errors.Is(err, ErrBadStrLen{"test", 2, 1}) {
		t.Error("bad minimum string length not detected")
	}
	if err.Error() != "field 'test' length 1 is less than 2" {
		t.Error("bad minimum string length error string")
	}

}
