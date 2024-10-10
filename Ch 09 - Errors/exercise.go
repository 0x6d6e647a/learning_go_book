package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Employee struct {
	ID        string
	FirstName string
	LastName  string
}

var ErrInvalidID = errors.New("invalid ID")

type EmptyFieldError struct {
	FieldName string
}

func (efe EmptyFieldError) Error() string {
	return efe.FieldName
}

var validID = regexp.MustCompile(`\w{4}-\d{3}`)

func ValidateEmployee(e Employee) error {
	var allErrors []error

	if len(e.ID) == 0 {
		allErrors = append(allErrors, EmptyFieldError{"ID"})
	}
	if !validID.MatchString(e.ID) {
		allErrors = append(allErrors, ErrInvalidID)
	}
	if len(e.FirstName) == 0 {
		allErrors = append(allErrors, EmptyFieldError{"FirstName"})
	}
	if len(e.LastName) == 0 {
		allErrors = append(allErrors, EmptyFieldError{"LastName"})
	}

	switch len(allErrors) {
	case 0:
		return nil
	case 1:
		return allErrors[0]
	default:
		return errors.Join(allErrors...)
	}
}

func errorString(err error, employee Employee, sb *strings.Builder) {
	var fieldErr EmptyFieldError
	if errors.Is(err, ErrInvalidID) {
		sb.WriteString(fmt.Sprintf("invalid ID: \"%s\"", employee.ID))
	} else if errors.As(err, &fieldErr) {
		sb.WriteString(fmt.Sprintf("empty field \"%s\"", fieldErr.FieldName))
	} else {
		sb.WriteString(fmt.Sprintf("%v", err))
	}
}

var Employees = [...]Employee{
	{"mndz-762", "Anthony", "Mendez"},
	{"blnk-000", "", ""},
	{"pike-123", "Rob", "Pike"},
	{"blnk-001", "Blank", ""},
	{"thom-321", "Ken", "Thompson"},
	{"blnk-002", "", "Blank"},
	{"bob-777", "Bob", "Bobberson"},
	{"", "", ""},
}

func main() {
	for index, employee := range Employees {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("record %d: %+v", index, employee))

		err := ValidateEmployee(employee)

		if err != nil {
			switch err := err.(type) {
			case interface{ Unwrap() []error }:
				sb.WriteString(" allErrors: ")
				for index, err := range err.Unwrap() {
					if index != 0 {
						sb.WriteString(", ")
					}
					errorString(err, employee, &sb)
				}
			default:
				sb.WriteString(" error: ")
				errorString(err, employee, &sb)
			}
		}
		fmt.Println(sb.String())
	}
}
