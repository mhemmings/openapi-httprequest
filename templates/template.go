package templates

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateArg struct {
	Types              DefinitionList
	Handlers           HandlerList
	GenerateServerCode bool
}

// WriteAll writes the generated packed to the provided outputDir
func WriteAll(outputDir string, args TemplateArg) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}
	err := Write(Params, args, filepath.Join(outputDir, "params.go"))
	if err != nil {
		return err
	}

	err = Write(Handlers, args, filepath.Join(outputDir, "handlers.go"))
	if err != nil {
		return err
	}

	if args.GenerateServerCode {
		err = Write(Main, args, filepath.Join(outputDir, "main.go"))
	}
	return err
}

func Write(template *template.Template, data TemplateArg, filepath string) error {
	// TODO: This is gross and broken. Find the proper way of adding required imports.
	var args = struct {
		TemplateArg
		Imports []string
	}{data, []string{}}

Outer:
	for _, def := range data.Types {
		if def.TypeStr == "time.Time" {
			args.Imports = append(args.Imports, "time")
			break Outer
		}
		for _, prop := range def.Properties {
			if prop.TypeStr == "time.Time" {
				args.Imports = append(args.Imports, "time")
				break Outer
			}
		}
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, args); err != nil {
		return err
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = dst.Write(source)
	return err
}
