package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"openapi-httprequest": main1,
	}))
}

func TestScripts(t *testing.T) {
	p := testscript.Params{
		Dir: "testdata",
	}
	err := gotooltest.Setup(&p)
	if err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}
