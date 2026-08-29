package palidate

import (
	"reflect"
	"regexp"
	"testing"
)

func ptr[T any](v T) *T {
	return &v
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		rawTag   string
		want     *Rule
		wantErrs int
	}{
		{
			name:     "empty tag",
			rawTag:   "",
			want:     &Rule{},
			wantErrs: 0,
		},
		{
			name:     "required only",
			rawTag:   "required",
			want:     &Rule{required: true},
			wantErrs: 0,
		},
		{
			name:     "minimum only",
			rawTag:   "min=3",
			want:     &Rule{min: ptr(3)},
			wantErrs: 0,
		},
		{
			name:     "maximum only",
			rawTag:   "max=20",
			want:     &Rule{max: ptr(20)},
			wantErrs: 0,
		},
		{
			name:     "regex only",
			rawTag:   "regex='^[a-z]+$'",
			want:     &Rule{pattern: regexp.MustCompile("^[a-z]+$")},
			wantErrs: 0,
		},
		{
			name:   "all validators",
			rawTag: "required,min=4,max=30,regex='^[A-Z]+$'",
			want: &Rule{
				required: true,
				min:      ptr(4),
				max:      ptr(30),
				pattern:  regexp.MustCompile("^[A-Z]+$"),
			},
			wantErrs: 0,
		},
		{
			name:   "validators in a different order",
			rawTag: "regex='hello .*',max=9,min=5,required",
			want: &Rule{
				required: true,
				min:      ptr(5),
				max:      ptr(9),
				pattern:  regexp.MustCompile("hello .*"),
			},
			wantErrs: 0,
		},
		{
			name:   "zero and negative bounds",
			rawTag: "min=-1,max=0",
			want: &Rule{
				min: ptr(-1),
				max: ptr(0),
			},
			wantErrs: 0,
		},
		{
			name:   "escape in regex",
			rawTag: "min=-10,regex='^[a-z'']+',max=999",
			want: &Rule{
				min:     ptr(-10),
				max:     ptr(999),
				pattern: regexp.MustCompile("^[a-z']+"),
			},
			wantErrs: 0,
		},
		{
			name:     "leading and trailing whitespace",
			rawTag:   "   required		",
			want:     &Rule{required: true},
			wantErrs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{}
			got, errs := p.Parse(tt.rawTag)
			if len(errs) != tt.wantErrs {
				t.Fatalf("parseTag: got %d errors, but wanted %d: %v", len(errs), tt.wantErrs, errs)
			}

			// unset the regex patterns since those can not be compared using reflect.DeepEqual.
			// Instead there is custom logic to compare those fields
			wantPat := tt.want.pattern
			tt.want.pattern = nil
			gotPat := got.pattern
			got.pattern = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTag: got '%v' but wanted '%v'", got, tt.want)
			}

			if (wantPat == nil) != (gotPat == nil) {
				t.Fatalf("parseTag: wanted '%v' but got '%v' instead", wantPat, gotPat)
			}

			if wantPat != nil && wantPat.String() != gotPat.String() {
				t.Errorf("parseTag: wanted pattern '%s' but got '%s'", wantPat.String(), gotPat.String())
			}
		})
	}
}
