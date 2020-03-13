package common

import "testing"

func TestNewSemver(t *testing.T) {
	type Expected struct {
		major, minor, patch int
	}

	type PassCase struct {
		version string
		expected Expected
	}

	passCases := []PassCase{
		{"0.17.3", Expected{0, 17, 3}},
		{"0.18", Expected{0, 18, 0}},
		{"32", Expected{32, 0, 0}},
	}

	for _, c := range passCases {
		actual, e := NewSemver(c.version)

		if e != nil {
			t.Errorf("NewSemver(%s): %s", c.version, e)
			continue
		}

		if !(actual.Major == c.expected.major && actual.Minor == c.expected.minor && actual.Patch == c.expected.patch) {
			t.Errorf("NewSemver(%s) = %d.%d.%d, expected: %d.%d.%d", c.version, actual.Major, actual.Minor, actual.Patch, c.expected.major, c.expected.minor, c.expected.patch)
		}
	}

	failCases := []string{"abcd", ""}

	for _, version := range failCases {
		actual, e := NewSemver(version)

		if e == nil {
			t.Errorf("NewSemver(%s) = %d.%d.%d, expected error", version, actual.Major, actual.Minor, actual.Patch)
		}
	}
}

func TestCmp(t *testing.T) {
	cases := []struct{
		a, b, expected int
	}{
		{0, 1, -1}, // a < b
		{1, 0, 1}, // a > b
		{3, 3, 0}, // a == b
	}

	for _, c := range cases {
		actual := cmp(c.a, c.b)

		if actual != c.expected {
			t.Errorf("cmp(%d, %d) = %d, expected: %d", c.a, c.b, actual, c.expected)
		}
	}
}

func TestSemverCmp(t *testing.T) {
	cases := []struct{
		a, b string
		expected int
	}{
		{"1.1.1", "2.2.2", -1},
		{"0.18.3", "0.18.3", 0},
		{"0.7", "0.2.3", 1},
	}

	for _, c := range cases {
		a, e := NewSemver(c.a)

		if e != nil {
			t.Fatal("TestSemverCmp:", e)
		}

		b, e := NewSemver(c.b)

		if e != nil {
			t.Fatal("TestSemverCmp:", e)
		}

		if actual := a.Cmp(b); actual != c.expected {
			t.Errorf("(%v).Cmp(%v) = %d, expected: %d", a, b, actual, c.expected)
		}
	}
}

func TestSemverCmpString(t *testing.T) {
	//
}
