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
	for f, fieldValue := range value.Fields() {
		tag, ok := f.Tag.Lookup("palidate")
		if !ok {
			// nothing to palidate here
			continue
		}

		p := parseTag(tag)

		if p.required && fieldValue.IsZero() {
			b.add(f.Name, "required field had it's zero value")
		}

		errMsg, ok := validMin(fieldValue, p.min)
		if !ok {
			b.add(f.Name, errMsg)
		}

		errMsg, ok = validMax(fieldValue, p.max)
		if !ok {
			b.add(f.Name, errMsg)
		}

		errMsg, ok = validPattern(fieldValue, p.pattern)
		if !ok {
			b.add(f.Name, errMsg)
		}
	}

	if b.hasErrors() {
		return b.Error()
	}

	return nil
}

func validMin(v reflect.Value, min *int) (string, bool) {
	if min == nil {
		// nothing to validate
		return "", true
	}

	if v.Kind() == reflect.String && v.Len() < *min {
		return "string was too short", false
	}

	if v.CanInt() && int(v.Int()) < *min {
		return "int was too small", false
	}

	return "", true
}

func validMax(v reflect.Value, max *int) (string, bool) {
	if max == nil {
		// nothing to validate
		return "", true
	}

	if v.Kind() == reflect.String && v.Len() > *max {
		return "string was too long", false
	}

	if v.CanInt() && int(v.Int()) > *max {
		return "int was too large", false
	}

	return "", true
}

func validPattern(v reflect.Value, pat *regexp.Regexp) (string, bool) {
	if pat == nil {
		// nothing to validate
		return "", true
	}

	if v.Kind() != reflect.String {
		return "patterns are only valid on string fields", false
	}

	if !pat.MatchString(v.String()) {
		return "failed to match string field against pattern", false
	}

	return "", true
}
