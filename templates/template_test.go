package templates

import "testing"

var commentTests = []struct {
	Input    string
	Expected string
}{{
	Input:    "Test comment",
	Expected: "// Test comment",
}, {
	Input:    "Test comment\nwhich is multiline",
	Expected: "// Test comment\n// which is multiline",
},
	{
		Input:    "",
		Expected: "",
	}}

func TestComment(t *testing.T) {
	for i, test := range commentTests {
		result := Comment(test.Input)
		if result != test.Expected {
			t.Errorf(`[%d] expected "%s", got "%s"`, i, test.Expected, result)
		}
	}
}
