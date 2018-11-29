package openapi

import (
	"testing"
)

var typesTestTable = []struct {
	Type     string
	Format   string
	Expected string
}{{
	Type:     "integer",
	Expected: "int64",
}, {
	Type:     "integer",
	Format:   "int32",
	Expected: "int",
}, {
	Type:     "integer",
	Format:   "int64",
	Expected: "int64",
}, {
	Type:     "number",
	Expected: "float64",
}, {
	Type:     "number",
	Format:   "float",
	Expected: "float64",
}, {
	Type:     "number",
	Format:   "double",
	Expected: "float64",
}, {
	Type:     "string",
	Expected: "string",
}, {
	Type:     "string",
	Format:   "byte",
	Expected: "string",
}, {
	Type:     "string",
	Format:   "binary",
	Expected: "string",
}, {
	Type:     "string",
	Format:   "date",
	Expected: "string",
}, {
	Type:     "string",
	Format:   "password",
	Expected: "string",
}, {
	Type:     "string",
	Format:   "date-time",
	Expected: "time.Time",
}, {
	Type:     "boolean",
	Expected: "bool",
}, {
	Type:     "unknown",
	Expected: "unknown",
}}

func TestOpenapiTypeToGo(t *testing.T) {
	for i, test := range typesTestTable {
		output := TypeString(test.Type, test.Format)
		if output != test.Expected {
			t.Errorf(`[%d] expected: "%s", got: "%s"`, i, test.Expected, output)
		}
	}
}
