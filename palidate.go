package palidate

import (
	"errors"
	"reflect"
	"regexp"
)

// Struct validates the struct using the palidate tags to determine how fields should be populated.
// v must be either a struct or a pointer to a struct, all validation failures are returned in an error.
// A nil return means that validation passed.
func Struct(v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Pointer {
		// indirect the pointer
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return errors.New("can not validate non-struct type")
	}

	b := ErrorBuilder{}
	p := Parser{}
	for f, fieldValue := range value.Fields() {
		rawTag, ok := f.Tag.Lookup("palidate")
		if !ok {
			// nothing to palidate here
			continue
		}

		tag, errs := p.Parse(rawTag)
		if errs != nil {
			b.addTagErrs(f.Name, errs...)
			continue
		}

		if tag.required && fieldValue.IsZero() {
			b.addFieldErrs(f.Name, errors.New("required field had it's zero value"))
		}

		err := validateMin(fieldValue, tag.min)
		if err != nil {
			b.addFieldErrs(f.Name, err)
		}

		err = validateMax(fieldValue, tag.max)
		if err != nil {
			b.addFieldErrs(f.Name, err)
		}

		err = validatePattern(fieldValue, tag.pattern)
		if err != nil {
			b.addFieldErrs(f.Name, err)
		}
	}

	if b.hasErrors() {
		return b.Error()
	}

	return nil
}

func validateMin(v reflect.Value, min *int) error {
	if min == nil {
		// nothing to validate
		return nil
	}

	if v.Kind() == reflect.String && v.Len() < *min {
		return errors.New("string was too short")
	}

	if v.CanInt() && int(v.Int()) < *min {
		return errors.New("int was too small")
	}

	return nil
}

func validateMax(v reflect.Value, max *int) error {
	if max == nil {
		// nothing to validate
		return nil
	}

	if v.Kind() == reflect.String && v.Len() > *max {
		return errors.New("string was too long")
	}

	if v.CanInt() && int(v.Int()) > *max {
		return errors.New("int was too large")
	}

	return nil
}

func validatePattern(v reflect.Value, pat *regexp.Regexp) error {
	if pat == nil {
		// nothing to validate
		return nil
	}

	if v.Kind() != reflect.String {
		return errors.New("patterns are only valid on string fields")
	}

	if !pat.MatchString(v.String()) {
		return errors.New("failed to match string field against pattern")
	}

	return nil
}
