package palidate

import "strings"

type FieldError struct {
	FieldName string
	Error     string
}

type ErrorBuilder struct {
	errors []FieldError
}

func (e *ErrorBuilder) add(fieldname, errMsg string) {
	e.errors = append(e.errors, FieldError{FieldName: fieldname, Error: errMsg})
}

func (e *ErrorBuilder) hasErrors() bool {
	return len(e.errors) != 0
}

func (e *ErrorBuilder) Error() error {
	return Error(e.errors)
}

type Error []FieldError

func (e Error) Error() string {
	b := strings.Builder{}
	b.WriteString("invalid struct fields| ")

	for i, f := range e {
		b.WriteString(f.FieldName)
		b.WriteString(": ")
		b.WriteString(f.Error)

		if i < len(e)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}
