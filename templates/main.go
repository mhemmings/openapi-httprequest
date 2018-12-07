package templates

import "text/template"

var Main = template.Must(template.New("").Parse(`
// The code in this file was automatically generated.
package main

import (
  "context"
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/julienschmidt/httprouter"
  httprequest "gopkg.in/httprequest.v1"
)

func main() {
  // TODO This should use the interface APIHandler
  f := func(p httprequest.Params) (*Handler, context.Context, error) {
    log.Printf("%s %s\n", p.Request.Method, p.Request.URL)
    return &Handler{}, p.Context, nil
  }

  router := httprouter.New()
  var reqSrv httprequest.Server
  for _, h := range reqSrv.Handlers(f) {
    router.Handle(h.Method, h.Path, h.Handle)
  }

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
`))
