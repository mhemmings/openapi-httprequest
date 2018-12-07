package openapi

import (
	"testing"
)

var requestTestTable = []struct {
	Location string
	Expected string
}{{
	Location: "path",
	Expected: "path",
}, {
	Location: "query",
	Expected: "form",
}, {
	Location: "header",
	Expected: "header",
}}

func TestParamLocation(t *testing.T) {
	for i, test := range requestTestTable {
		output := ParamLocation(test.Location)
		if output != test.Expected {
			t.Errorf(`[%d] expected: "%s", got: "%s"`, i, test.Expected, output)
		}
	}
}
