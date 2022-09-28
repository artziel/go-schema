package schema

import (
	"testing"
)

func TestToPair(t *testing.T) {

	type Test struct {
		Input    string
		Expected pair
	}

	tests := []Test{
		{
			Input:    "TheName:TheValue",
			Expected: pair{Name: "thename", Value: "TheValue"},
		},
		{
			Input:    "The Name:The:Value",
			Expected: pair{Name: "the name", Value: "The:Value"},
		},
		{
			Input:    "The Name",
			Expected: pair{Name: "the name", Value: ""},
		},
	}

	for i, test := range tests {
		p := toPair(test.Input, ':')
		if p.Name != test.Expected.Name || p.Value != test.Expected.Value {
			t.Errorf("Test %d Fail, result dont match expected value:\nGot  %v\nWant %v", i, p, test.Expected)
		}
	}

}

func TestSplitAt(t *testing.T) {

	type Test struct {
		Input    string
		Expected []string
	}

	tests := []Test{
		{
			Input:    "this is a test, split expected in N parts",
			Expected: []string{"this", "is", "a", "test,", "split", "expected", "in", "N", "parts"},
		},
		{
			Input:    "the quote string \"This is a Quote String\" sould not be split",
			Expected: []string{"the", "quote", "string", "\"This is a Quote String\"", "sould", "not", "be", "split"},
		},
	}

	for i, test := range tests {
		p := splitAt(test.Input, ' ', '"')
		if len(p) != len(test.Expected) {
			t.Errorf("Test %d Fail, result dont match expected value:\nGot  %v\nWant %v", i, p, test.Expected)
		}
	}

}
