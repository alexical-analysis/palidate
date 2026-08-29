package palidate

import (
	"errors"
)

var (
	missingRequiredFieldErr = errors.New("required field had it's zero value")
	stringTooShortErr       = errors.New("string was too short")
	intTooSmallErr          = errors.New("int too small")
	stringTooLongErr        = errors.New("string was too long")
	intTooLargeErr          = errors.New("int was too large")
	patternForNonStringErr  = errors.New("patterns are only valid on string fields")
	patternMismatchErr      = errors.New("failed to match string field against pattern")
	malformedTag            = errors.New("failed to get palidate tag")
)

// Struct validates the struct using the palidate tags to determine how fields should be populated.
// v must be either a struct or a pointer to a struct, all validation failures are returned in an error.
// A nil return means that validation passed.
func Struct(v any) error {
	//
	// TODO: make sure v is either a struct or struct pointer
	//

	// TODO: use the error builder to build a clean set of errors for all the struct fields
	b := ErrorBuilder{}
	b.addFieldErrs("field", missingRequiredFieldErr)

	// TODO: use the parser to parse the palidate tags into palidate.Rule structs for easy use
	p := Parser{}
	_, errs := p.Parse("valid")
	if errs != nil {
		b.addTagErrs("field", errs...)
	}

	//
	// TODO: loop through fields validating them against the palidate.Rule
	//

	if b.hasErrors() {
		return b.Error()
	}

	return nil
}
