package helper

import (
	"fmt"
	"testing"
)

func TestCleanBytes(t *testing.T) {

	a := string([]byte{8, 9, 10, 12, 13})
	//a = string([]byte{10})

	fmt.Println(a)
}

func TestSubString(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		start  int
		length int
		want   string
	}{
		{"normal", "hello world", 0, 5, "hello"},
		{"out of bounds start", "hello", 10, 5, ""},
		{"out of bounds length", "hello", 0, 10, "hello"},
		{"negative start", "hello", -1, 5, "hello"},
		{"negative length", "hello", 0, -1, ""},
		{"partial overlap", "hello", 3, 5, "lo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubString(tt.src, tt.start, tt.length); got != tt.want {
				t.Errorf("SubString() = %v, want %v", got, tt.want)
			}
		})
	}
}
