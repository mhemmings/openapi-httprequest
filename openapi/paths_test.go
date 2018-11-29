package openapi

import (
	"testing"
)

var pathsTestTable = []struct {
	Path     string
	Expected string
}{{
	Path:     "/test/{foo}",
	Expected: "/test/:foo",
}, {
	Path:     "/test/foo",
	Expected: "/test/foo",
}, {
	Path:     "/test/{foo}/{bar}",
	Expected: "/test/:foo/:bar",
}}

func TestOpenApiPath(t *testing.T) {
	for i, test := range pathsTestTable {
		output := PathToString(test.Path)
		if output != test.Expected {
			t.Errorf(`[%d] expected: "%s", got: "%s"`, i, test.Expected, output)
		}
	}
}
