package palidate

import (
	"testing"
)

func TestStructMinimum(t *testing.T) {
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

func TestStructExtended(t *testing.T) {
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
