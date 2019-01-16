package templates

import "text/template"

var GoMod = template.Must(template.New("").Parse(`
module openapi-httprequest-server

require (
	gopkg.in/httprequest.v1 v1.1.3
)
`))
