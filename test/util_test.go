package test

import (
	"testing"

	"github.com/tkeyo/tinyml-be/util"
)

func TestIsAuthorized(t *testing.T) {
	type TestCase struct {
		input    string
		expected bool
	}

	cases := []TestCase{
		{"000", false},
	}

	for _, c := range cases {
		actual := util.IsAuthorized(c.input)
		expected := c.expected

		if actual != expected {
			t.Fatalf("Expected %t, got %t", expected, actual)
		}
	}
}
