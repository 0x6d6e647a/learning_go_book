package ch16

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var ErrNonStruct = errors.New("not a struct")

type ErrStrConv struct {
	FieldName string
	MinStrLen string
	Err       error
}

func (e ErrStrConv) Error() string {
	return fmt.Sprintf(
		"field '%s' minStrLen tag '%s' is not an int",
		e.FieldName,
		e.MinStrLen,
	)
}

type ErrBadStrLen struct {
	FieldName string
	MinStrLen int
	StrLen    int
}

func (e ErrBadStrLen) Error() string {
	return fmt.Sprintf(
		"field '%s' length %d is less than %d",
		e.FieldName,
		e.StrLen,
		e.MinStrLen,
	)
}

func ValidateStringLength(v any) (errs error) {
	vt := reflect.TypeOf(v)

	// Invalid input.
	if vt.Kind() != reflect.Struct {
		errs = errors.Join(errs, ErrNonStruct)
		return errs
	}

	// Check each field in structure.
	vv := reflect.ValueOf(v)

	for i := 0; i < vt.NumField(); i += 1 {
		field := vt.Field(i)

		// Skip non-string fields.
		if field.Type.Kind() != reflect.String {
			continue
		}

		// Skip fields without 'minStrLen' tagStr.
		tagStr, ok := field.Tag.Lookup("minStrLen")
		if !ok {
			continue
		}

		// Convert tag to int.
		minStrLen, err := strconv.Atoi(tagStr)
		if err != nil {
			errs = errors.Join(errs, ErrStrConv{
				field.Name,
				tagStr,
				err,
			})
		}

		// Check length of string.
		strLen := len(vv.Field(i).String())
		if strLen < minStrLen {
			errs = errors.Join(errs, ErrBadStrLen{
				field.Name,
				minStrLen,
				strLen,
			})
		}
	}

	return errs
}
