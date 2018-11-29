package openapi

import (
  "io/ioutil"
  "net/http"
  "strings"

  "github.com/getkin/kin-openapi/openapi3"
)

// Load takes either a a local file path or a http path and loads it into a openapi3.Swagger.
func Load(uri string) (*openapi3.Swagger, error) {
  var data []byte
  var err error

  // Load either from an HTTP URL or from a local file depending on the passed value.
  if strings.HasPrefix(uri, "http") {
    resp, err := http.Get(uri)
    if err != nil {
      return nil, err
    }

    data, err = ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
      return nil, err
    }
  } else {
    data, err = ioutil.ReadFile(uri)
    if err != nil {
      return nil, err
    }
  }

  // Load the OpenAPI document.
  loader := openapi3.NewSwaggerLoader()
  if strings.HasSuffix(uri, ".yaml") || strings.HasSuffix(uri, ".yml") {
    return loader.LoadSwaggerFromYAMLData(data)
  }

  return loader.LoadSwaggerFromData(data)
}
