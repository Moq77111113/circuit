package tags

import (
	"reflect"
	"testing"
	"time"
)

func TestDereferenceType(t *testing.T) {
	tests := []struct {
		name string
		in   reflect.Type
		want reflect.Type
	}{
		{
			name: "string",
			in:   reflect.TypeOf(""),
			want: reflect.TypeOf(""),
		},
		{
			name: "pointer to string",
			in:   reflect.TypeOf(new(string)),
			want: reflect.TypeOf(""),
		},
		{
			name: "double pointer to int",
			in:   reflect.TypeOf(new(*int)),
			want: reflect.TypeOf(0),
		},
		{
			name: "pointer to struct",
			in:   reflect.TypeOf(&time.Time{}),
			want: reflect.TypeOf(time.Time{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dereferenceType(tt.in)
			if got != tt.want {
				t.Errorf("dereferenceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElementType(t *testing.T) {
	tests := []struct {
		name        string
		in          reflect.Type
		wantType    reflect.Type
		wantIsSlice bool
	}{
		{
			name:        "slice of string",
			in:          reflect.TypeOf([]string{}),
			wantType:    reflect.TypeOf(""),
			wantIsSlice: true,
		},
		{
			name:        "slice of int",
			in:          reflect.TypeOf([]int{}),
			wantType:    reflect.TypeOf(0),
			wantIsSlice: true,
		},
		{
			name:        "slice of struct",
			in:          reflect.TypeOf([]time.Time{}),
			wantType:    reflect.TypeOf(time.Time{}),
			wantIsSlice: true,
		},
		{
			name:        "not a slice",
			in:          reflect.TypeOf(""),
			wantType:    reflect.TypeOf(""),
			wantIsSlice: false,
		},
		{
			name:        "array (not slice)",
			in:          reflect.TypeOf([3]int{}),
			wantType:    reflect.TypeOf([3]int{}),
			wantIsSlice: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotIsSlice := elementType(tt.in)
			if gotType != tt.wantType {
				t.Errorf("elementType() type = %v, want %v", gotType, tt.wantType)
			}
			if gotIsSlice != tt.wantIsSlice {
				t.Errorf("elementType() isSlice = %v, want %v", gotIsSlice, tt.wantIsSlice)
			}
		})
	}
}
