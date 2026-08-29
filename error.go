package palidate

import (
	"errors"
	"strings"
)

// ErrorBuilder can be used to build a palidate Error with zero or more tag and field errors
type ErrorBuilder struct {
	fieldErrors map[string][]error
	tagErrors   map[string][]error
}

// addFieldErrs adds zero or more field errors for the given fieldName.
func (e *ErrorBuilder) addFieldErrs(fieldName string, err ...error) {
	if e.fieldErrors == nil {
		e.fieldErrors = make(map[string][]error)
	}

	e.fieldErrors[fieldName] = append(e.fieldErrors[fieldName], err...)
}

// addTagErrs adds zero or more tag errors for the given fieldName.
func (e *ErrorBuilder) addTagErrs(fieldName string, err ...error) {
	if e.tagErrors == nil {
		e.tagErrors = make(map[string][]error)
	}

	e.tagErrors[fieldName] = append(e.tagErrors[fieldName], err...)
}

// hasErrors returns true if any field has at least one field or tag error
func (e *ErrorBuilder) hasErrors() bool {
	return len(e.fieldErrors) > 0 || len(e.tagErrors) > 0
}

// Error builds the actual error from the provided field and tag errors
func (e *ErrorBuilder) Error() error {
	tmpMap := make(map[string][]error)

	for field, errs := range e.fieldErrors {
		tmpMap[field] = append(tmpMap[field], errs...)
	}

	for field, errs := range e.tagErrors {
		tmpMap[field] = append(tmpMap[field], errs...)
	}

	errMap := make(map[string]error, len(tmpMap))
	for field, errs := range tmpMap {
		errMap[field] = joinErros(errs)
	}

	return Error{errMap: errMap}
}

func joinErros(errs []error) error {
	s := []string{}
	for _, err := range errs {
		if err != nil {
			s = append(s, err.Error())
		}
	}
	if len(s) == 0 {
		return nil
	}

	return errors.New(strings.Join(s, ", "))
}

// Error contains all the errors for all the struct fields validated by palidate
type Error struct {
	errMap map[string]error
}

// Error implements the error interface for the Error type
func (e Error) Error() string {
	b := strings.Builder{}
	if len(e.errMap) > 0 {
		b.WriteString("invalid struct")
	}

	for field, err := range e.errMap {
		b.WriteString("| ")
		b.WriteString(field)
		b.WriteString(": ")
		b.WriteString(err.Error())
	}

	return b.String()
}
