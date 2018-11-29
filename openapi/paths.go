package openapi

import "strings"

// PathToString translates an OpenAPI path (/{foo} style) into httprequest path (/:foo style)
func PathToString(path string) string {
	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)
	return path
}
