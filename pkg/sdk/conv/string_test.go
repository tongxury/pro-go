package conv

import "testing"

func TestStr(t *testing.T) {

	cases := []struct {
		in   any
		want string
	}{
		{in: "1", want: "1"},
		{in: 1, want: "1"},
		{in: int64(1), want: "1"},
		{in: int32(1), want: "1"},
		{in: int8(1), want: "1"},
		{in: uint8(1), want: "1"},
		{in: float64(1), want: "1"},
	}

	for _, c := range cases {

		got := Str(c.in)
		if got != c.want {
			t.Errorf("Str(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}
