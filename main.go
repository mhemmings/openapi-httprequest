package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/juju/gnuflag"
	oas "github.com/mhemmings/openapi-httprequest/openapi"
	"github.com/mhemmings/openapi-httprequest/templates"
	errgo "gopkg.in/errgo.v1"
)

type ref struct {
	Name      string
	SchemaRef *openapi3.SchemaRef
}

var references = make(map[string]ref)

var printCmdUsage = func() {
	fmt.Printf("usage: openapi-httprequest [flags] openapidoc.yaml\n\n")
	gnuflag.PrintDefaults()
}

var (
	outputDir      = gnuflag.String("outputdir", "", "The output directory to save generated server package (default: the current directory, or a temporary directory if --http is specified")
	listenAddr     = gnuflag.String("http", "", "Implies --server. If set, the generated server will be run on the given network address (e.g. localhost:8088)")
	packageName    = gnuflag.String("pkg", "params", "Package name to use for generated files (ignored if --server is specified)")
	generateServer = gnuflag.Bool("server", false, "Generate server code (overwrites --pkg=main)")
)

func main() {
	os.Exit(main1())
}

func main1() int {
	gnuflag.Usage = func() {
		printCmdUsage()
	}
	log.SetFlags(0)

	gnuflag.Parse(true)
	if gnuflag.NArg() != 1 || gnuflag.Arg(0) == "help" {
		gnuflag.Usage()
		return 2
	}
	if err := main2(); err != nil {
		log.Printf("%v", err)
		return 1
	}
	return 0
}

func main2() error {
	if *outputDir == "" {
		if *listenAddr != "" {
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				return err
			}
			defer os.RemoveAll(dir)
			*outputDir = dir
		} else {
			*outputDir = "."
		}
	}
	if *listenAddr != "" {
		*generateServer = true
	}
	if *generateServer {
		*packageName = "main"
	}
	uri := gnuflag.Arg(0)
	swagger, err := oas.Load(uri)
	if err != nil {
		return errgo.Mask(err)
	}

	arg := templates.TemplateArg{
		GenerateServer: *generateServer,
		Pkg:            *packageName,
	}

	// Build references of top level schema definitions
	for schemaName, schema := range swagger.Components.Schemas {
		references["#/components/schemas/"+schemaName] = ref{
			Name:      strcase.ToCamel(schemaName),
			SchemaRef: schema,
		}
	}

	// Build schema types
	for schemaName, schema := range swagger.Components.Schemas {
		s := schemaRefParse(schema, strcase.ToCamel(schemaName))
		arg.Types = append(arg.Types, &s)
	}

	// Sort schemas types so they appear in alphabetical order at the top of the file
	sort.Sort(arg.Types)

	// Build all the types for paths
	var reqResp templates.DefinitionList
	for path, pathItem := range swagger.Paths {
		for method, op := range pathItem.Operations() {
			if method == "HEAD" || method == "OPTIONS" {
				// Ignore (https://github.com/go-httprequest/httprequest/blob/2b21a94c9e788981d4e609ef4b7a21cedae6da66/type.go#L225)
				continue
			}
			if op.OperationID == "" {
				return errgo.Newf("missing operationId for path %s %s", method, path)
			}
			name := strcase.ToCamel(op.OperationID + "Request")
			req := templates.Definition{
				Name:       name,
				DocComment: templates.Comment(fmt.Sprintf("%s holds the request value for a %s call.", name, op.OperationID)),
				// Embed the the httprequest.Route type
				Properties: templates.DefinitionList{{
					Tag:     fmt.Sprintf("`httprequest:\"%s %s\"`", method, oas.PathToString(path)),
					TypeStr: "httprequest.Route",
				}},
			}

			if op.Description != "" {
				// If this request has a description, also include it in the comment.
				req.DocComment = req.DocComment + "\n" + templates.Comment(op.Description)
			}

			handler := templates.Handler{
				Request: req.Name,
			}

			// Get request params
			for _, param := range op.Parameters {
				def := schemaRefParse(param.Value.Schema, strcase.ToCamel(param.Value.Name))
				def.TypeStr = strings.TrimPrefix(def.TypeStr, "*")
				p := templates.Definition{
					Name:    def.Name,
					Tag:     fmt.Sprintf("`httprequest:\"%s,%s\"`", param.Value.Name, oas.ParamLocation(param.Value.In)),
					TypeStr: def.TypeStr,
				}

				req.Properties = append(req.Properties, &p)
			}

			// Get request body
			if op.RequestBody != nil && op.RequestBody.Value.Content["application/json"] != nil {
				if schema := op.RequestBody.Value.Content["application/json"].Schema; schema != nil {
					def := schemaRefParse(schema, "")
					def.TypeStr = strings.TrimPrefix(def.TypeStr, "*")
					p := templates.Definition{
						Name:    "Body",
						Tag:     "`httprequest:\",body\"`",
						TypeStr: def.Name,
					}
					if def.Name == "" {
						// If the request body is not a referenced type and is instead defined inline,
						// we need to build the request type
						reqBody := def
						reqBody.Name = strcase.ToCamel(op.OperationID + "RequestBody")
						p.TypeStr = reqBody.Name
						reqResp = append(reqResp, &reqBody)
					}

					req.Properties = append(req.Properties, &p)
				}
			}
			reqResp = append(reqResp, &req)

			// Take the first response that isn't "default" and is a 2xx.
			// TODO: This needs much improvement.
			// TODO: Error if there's more than one 2xx response?
			for respName, response := range op.Responses {
				handler := handler
				if respName == "default" || !strings.HasPrefix(respName, "2") {
					// Don't build the "default" response as this is usually an error.
					// May not be the correct assumption.
					continue
				}

				name := op.OperationID
				var resp templates.Definition
				// If there's only one possible content type and it's JSON, then
				// we can use the httprequest method return type; otherwise
				// we'll leave the response empty and leave it up to the server
				// method to write the appropriate response.
				if len(response.Value.Content) == 1 {
					if body := response.Value.Content.Get("application/json"); body != nil {
						resp = schemaRefParse(body.Schema, "")
						// If the response is not a referenced type and is instead defined inline,
						// we need to build the response type.
						if resp.Name == "" {
							resp.Name = strcase.ToCamel(name + "Response")
							resp.DocComment = templates.Comment(fmt.Sprintf("%s holds the response to a %s call.", resp.Name, name))
							reqResp = append(reqResp, &resp)
						}
						handler.Response = resp.Name
					}
				}
				handler.Name = strcase.ToCamel(name)
				handler.DocComment = templates.Comment(op.Summary)
				arg.Handlers = append(arg.Handlers, &handler)
				break
			}
		}
	}

	sort.Sort(reqResp)
	arg.Types = append(arg.Types, reqResp...)

	sort.Sort(arg.Handlers)

	err = templates.WriteAll(*outputDir, arg)
	if err != nil {
		return errgo.Notef(err, "cannot write templates")
	}

	if *listenAddr != "" {
		fmt.Printf("Running API server at %s\n", *listenAddr)
		cmd := exec.Command("go", "build", "-o", "openapi-httprequest-server")
		cmd.Dir = *outputDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return errgo.Mask(err)
		}
		cmd = exec.Command(filepath.Join(*outputDir, "openapi-httprequest-server"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("LISTEN_ADDR=%s", *listenAddr))
		if err := cmd.Run(); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

// schemaRefParse takes an openapi SchemeRef doc and creates a type Definition to be used in params.go.
// It attempts ro recursively resolve all references.
func schemaRefParse(oasSchema *openapi3.SchemaRef, name string) templates.Definition {
	if oasSchema.Ref != "" {
		r := references[oasSchema.Ref]
		// If the definition name is currently unknown, set it to the resolved one.
		if name == "" || name == "interface{}" {
			name = r.Name
		}
		def := schemaRefParse(r.SchemaRef, name)
		def.TypeStr = r.Name
		if r.SchemaRef.Value.Type == "object" {
			def.TypeStr = "*" + def.TypeStr
		}
		def.DocComment = "" // Avoid duplication of comments due to references.
		return def
	}

	schema := templates.Definition{
		Name:       name,
		DocComment: templates.Comment(oasSchema.Value.Description),
	}

	if len(oasSchema.Value.Properties) > 0 {
		required := make(map[string]bool)
		for _, name := range oasSchema.Value.Required {
			required[name] = true
		}
		for propName, prop := range oasSchema.Value.Properties {
			p := schemaRefParse(prop, strcase.ToCamel(propName))
			omitempty := ",omitempty"
			if p.TypeStr == "" {
				p.TypeStr = "*" + p.Name
			}
			if required[propName] {
				p.TypeStr = strings.TrimPrefix(p.TypeStr, "*")
				omitempty = ""
			}
			p.Tag = fmt.Sprintf("`json:\"%s%s\"`", propName, omitempty)
			schema.Properties = append(schema.Properties, &p)
		}
		sort.Sort(schema.Properties)
	} else if oasSchema.Value.Items != nil {
		name := oas.TypeString(oasSchema.Value.Items.Value.Type, oasSchema.Value.Items.Value.Format)
		name = strings.TrimPrefix(name, "*")
		t := schemaRefParse(oasSchema.Value.Items, name)
		schema.TypeStr = fmt.Sprintf("[]%s", t.Name)
	} else { //native type
		schema.TypeStr = oas.TypeString(oasSchema.Value.Type, oasSchema.Value.Format)
	}

	return schema
}
