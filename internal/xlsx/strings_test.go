package xlsx

import "testing"

func TestNormalizeString(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"empty", "", ""},
		{"plain", "hello", "hello"},
		{"leading_trailing_space", "  hello  ", "hello"},
		{"tab_inside", "he\tllo", "hello"},
		{"newline_inside", "he\nllo", "hello"},
		{"ascii_space_preserved", "hello world", "hello world"},
		{"nbsp_removed", "hello world", "helloworld"},
		{"ideographic_space_removed", "hello　world", "helloworld"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeString(c.in)
			if got != c.want {
				t.Errorf("NormalizeString(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{" ", true},
		{"a", false},
		{" a ", false},
	}
	for _, c := range cases {
		if got := IsBlank(c.in); got != c.want {
			t.Errorf("IsBlank(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
