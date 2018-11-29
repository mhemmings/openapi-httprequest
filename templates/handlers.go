package templates

import "text/template"

// Handler defines a httprequest handler function to be written to handlers.go
type Handler struct {
  Name     string
  Request  string
  Response string
}

type HandlerList []*Handler

func (s HandlerList) Len() int           { return len(s) }
func (s HandlerList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s HandlerList) Less(i, j int) bool { return s[i].Name < s[j].Name }

var Handlers = template.Must(template.New("").Parse(`
// The code in this file was automatically generated.
package main

type Handler struct {
}
{{range .Handlers}}
func (Handler) {{.Name}}(req *{{.Request}}) ({{.Response}}, error) {
  resp := {{.Response}}{}
  return resp, nil
}
{{end}}`))
