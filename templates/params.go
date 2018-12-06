package templates

import "text/template"

// Definition defines a type definition to be written to params.go
type Definition struct {
	Name       string
	TypeStr    string
	Tag        string
	Properties DefinitionList
}

type DefinitionList []*Definition

func (s DefinitionList) Len() int           { return len(s) }
func (s DefinitionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s DefinitionList) Less(i, j int) bool { return s[i].Name < s[j].Name }

var Params = template.Must(template.New("").Parse(`
// The code in this file was automatically generated.
package main

import(
    {{range .Imports}}"{{.}}"{{end}}

    httprequest "gopkg.in/httprequest.v1"
)

{{range .Types}}
type {{.Name}} {{if .TypeStr}}{{.TypeStr}}{{else}}struct {
  {{range .Properties}}
  {{- .Name}} {{.TypeStr}} {{.Tag}}
  {{end}}
}{{end}}
{{end}}`))
