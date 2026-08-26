package palidate

import (
	"reflect"
	"regexp"
	"testing"
)

func ptr[T any](v T) *T {
	return &v
}

func Test_parseTag(t *testing.T) {
	tests := []struct {
		name   string
		rawTag string
		want   parsedTag
	}{
		{
			name:   "empty tag",
			rawTag: "",
			want:   parsedTag{},
		},
		{
			name:   "required only",
			rawTag: "required",
			want:   parsedTag{required: true},
		},
		{
			name:   "minimum only",
			rawTag: "min=3",
			want:   parsedTag{min: ptr[int](3)},
		},
		{
			name:   "maximum only",
			rawTag: "max=20",
			want:   parsedTag{max: ptr[int](20)},
		},
		{
			name:   "regex only",
			rawTag: "regex='^[a-z]+$'",
			want:   parsedTag{pattern: regexp.MustCompile("^[a-z]+$")},
		},
		{
			name:   "all validators",
			rawTag: "required,min=4,max=30,regex='^[A-Z]+$'",
			want: parsedTag{
				required: true,
				min:      ptr[int](4),
				max:      ptr[int](30),
				pattern:  regexp.MustCompile("^[A-Z]+$"),
			},
		},
		{
			name:   "validators in a different order",
			rawTag: "regex='hello .*',max=9,min=5,required",
			want: parsedTag{
				required: true,
				min:      ptr[int](5),
				max:      ptr[int](9),
				pattern:  regexp.MustCompile("hello .*"),
			},
		},
		{
			name:   "zero and negative bounds",
			rawTag: "min=-1,max=0",
			want: parsedTag{
				min: ptr[int](-1),
				max: ptr[int](0),
			},
		},
		{
			name:   "escape in regex",
			rawTag: "min=-10,regex='^[a-z'']+',max=999",
			want: parsedTag{
				min:     ptr[int](-10),
				max:     ptr[int](999),
				pattern: regexp.MustCompile("^[a-z']+"),
			},
		},
		{
			name:   "leading and trailing whitespace",
			rawTag: "   required		",
			want:   parsedTag{required: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTag(tt.rawTag)

			// unset the regex patterns since those can not be compared using reflect. Instead there
			// is custom logic to compare those fields
			wantPat := tt.want.pattern
			tt.want.pattern = nil
			gotPat := got.pattern
			got.pattern = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTag: got %v but wanted %v", got, tt.want)
			}

			if wantPat == nil && gotPat != nil {
				t.Fatalf("parseTag: wanted nil pattern but got a non-nil one")
			}

			if wantPat != nil && gotPat == nil {
				t.Fatalf("parseTag: wanted non-nil pattern but got a nil one")
			}

			if wantPat == nil {
				// nothing more to check
				return
			}

			if wantPat.String() != gotPat.String() {
				t.Errorf("parseTag: wanted pattern %s but got %s", wantPat.String(), gotPat.String())
			}
		})
	}
}
