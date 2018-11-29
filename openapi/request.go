package openapi

// ParamLocation translates an OpenAPI paramater location and translates it to a httprequest parameter location
func ParamLocation(loc string) string {
  switch loc {
  case "path":
    return "path"
  case "query":
    return "form"
  case "header":
    return "header"
  }
  return loc
}
