package palidate

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorBuilderAddFieldErrs(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want int
	}{
		{name: "no errors", want: 0},
		{name: "one error", errs: []error{errors.New("first")}, want: 1},
		{name: "multiple errors", errs: []error{errors.New("first"), errors.New("second")}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ErrorBuilder{}
			builder.addFieldErrs("field", tt.errs...)
			if got := len(builder.fieldErrors["field"]); got != tt.want {
				t.Errorf("addFieldErrs stored %d errors, wanted %d", got, tt.want)
			}
		})
	}
}

func TestErrorBuilderAddTagErrs(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want int
	}{
		{name: "no errors", want: 0},
		{name: "one error", errs: []error{errors.New("first")}, want: 1},
		{name: "multiple errors", errs: []error{errors.New("first"), errors.New("second")}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ErrorBuilder{}
			builder.addTagErrs("field", tt.errs...)
			if got := len(builder.tagErrors["field"]); got != tt.want {
				t.Errorf("addTagErrs stored %d errors, wanted %d", got, tt.want)
			}
		})
	}
}

func TestErrorBuilderHasErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*ErrorBuilder)
		want  bool
	}{
		{name: "empty", want: false},
		{
			name: "field errors",
			setup: func(builder *ErrorBuilder) {
				builder.addFieldErrs("field", errors.New("field error"))
			},
			want: true,
		},
		{
			name: "tag errors",
			setup: func(builder *ErrorBuilder) {
				builder.addTagErrs("field", errors.New("tag error"))
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ErrorBuilder{}
			if tt.setup != nil {
				tt.setup(&builder)
			}
			if got := builder.hasErrors(); got != tt.want {
				t.Errorf("hasErrors returned %t, wanted %t", got, tt.want)
			}
		})
	}
}

func TestErrorBuilderError(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*ErrorBuilder)
		wantMessage []string
	}{
		{
			name: "empty",
		},
		{
			name: "field and tag errors",
			setup: func(builder *ErrorBuilder) {
				builder.addFieldErrs("field", errors.New("field error"))
				builder.addTagErrs("tag", errors.New("tag error"))
			},
			wantMessage: []string{"field error", "tag error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := ErrorBuilder{}
			if tt.setup != nil {
				tt.setup(&builder)
			}

			got := builder.Error()
			for _, message := range tt.wantMessage {
				if !strings.Contains(got.Error(), message) {
					t.Errorf("Error did not contain %q: %q", message, got.Error())
				}
			}
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "empty",
			err:  Error{},
			want: "",
		},
		{
			name: "field error",
			err: Error{
				errMap: map[string]error{"field": errors.New("invalid value")},
			},
			want: "invalid struct| field: invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error returned %q, wanted %q", got, tt.want)
			}
		})
	}
}
