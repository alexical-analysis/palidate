package palidate

import (
	"testing"
)

func TestStructRequired(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "struct with no tags",
			v: struct {
				A int
				B int
			}{},
			want: "",
		},
		{
			name: "struct missing required field",
			v: struct {
				A string `palidate:"required"`
				B int
			}{},
			want: "invalid struct| A: required field had it's zero value",
		},
		{
			name: "struct has required field",
			v: struct {
				A string `palidate:"required"`
				B int
			}{A: "hello"},
			want: "",
		},
		{
			name: "pointer to struct",
			v: &struct {
				A string `palidate:"required"`
			}{A: "hello"},
			want: "",
		},
		{
			name: "non-struct value",
			v:    "hello",
			want: "can not validate non-struct type",
		},
		{
			name: "nil value",
			v:    nil,
			want: "can not validate non-struct type",
		},
		{
			name: "malformed tag",
			v: struct {
				A string `palidate:"unknown"`
			}{},
			want: "invalid struct| A: unknown rule 'unknown'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Struct(tt.v)
			errStr := ""
			if got != nil {
				errStr = got.Error()
			}

			if errStr != tt.want {
				t.Errorf("Struct: got '%v' but wanted '%v'", got, tt.want)
			}
		})
	}
}

func TestStructMaxMin(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "string field too short",
			v: struct {
				Phone string `palidate:"min=7,max=11"`
			}{
				Phone: "123",
			},
			want: "invalid struct| Phone: string was too short",
		},
		{
			name: "string field too long",
			v: struct {
				Phone string `palidate:"max=11"`
			}{
				Phone: "123456789012",
			},
			want: "invalid struct| Phone: string was too long",
		},
		{
			name: "valid string field",
			v: struct {
				Phone string `palidate:"min=3,max=11"`
			}{
				Phone: "hello",
			},
			want: "",
		},
		{
			name: "integer field too small",
			v: struct {
				Age int `palidate:"min=18"`
			}{
				Age: 17,
			},
			want: "invalid struct| Age: int was too small",
		},
		{
			name: "integer field too large",
			v: struct {
				Age int `palidate:"max=65"`
			}{
				Age: 66,
			},
			want: "invalid struct| Age: int was too large",
		},
		{
			name: "values within bounds",
			v: struct {
				Name string `palidate:"min=3,max=10"`
				Age  int    `palidate:"min=18,max=65"`
			}{
				Name: "Alex",
				Age:  30,
			},
			want: "",
		},
		{
			name: "non-string and non-integer field",
			v: struct {
				Enabled bool `palidate:"min=1,max=0"`
			}{
				Enabled: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Struct(tt.v)
			errStr := ""
			if got != nil {
				errStr = got.Error()
			}

			if errStr != tt.want {
				t.Errorf("Struct: got '%v' but wanted '%v'", got, tt.want)
			}
		})
	}
}

func TestStructRegex(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "matching string field",
			v: struct {
				Name string `palidate:"regex='^[a-z]+$'"`
			}{
				Name: "alex",
			},
			want: "",
		},
		{
			name: "non-matching string field",
			v: struct {
				Name string `palidate:"regex='^[a-z]+$'"`
			}{
				Name: "Alex123",
			},
			want: "invalid struct| Name: failed to match string field against pattern",
		},
		{
			name: "non-string field",
			v: struct {
				Age int `palidate:"regex='^[0-9]+$'"`
			}{
				Age: 30,
			},
			want: "invalid struct| Age: patterns are only valid on string fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Struct(tt.v)
			errStr := ""
			if got != nil {
				errStr = got.Error()
			}

			if errStr != tt.want {
				t.Errorf("Struct: got '%v' but wanted '%v'", got, tt.want)
			}
		})
	}
}
