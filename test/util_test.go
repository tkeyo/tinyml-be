package test

import (
	"testing"

	"github.com/tkeyo/tinyml-be/util"
)

// Test authorization
func TestIsAuthorized(t *testing.T) {
	type TestCase struct {
		requestAuthKey string
		APIAuthKey     string
		expected       bool
	}

	cases := []TestCase{
		{"000", "123", false},
		{"123", "123", true},
	}

	for _, c := range cases {
		actual := util.IsAuthorized(c.requestAuthKey, c.APIAuthKey)
		expected := c.expected

		if actual != expected {
			t.Fatalf("Expected %t, got %t", expected, actual)
		}
	}
}
