package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/thema"
	"github.com/grafana/thema/load"
	"github.com/grafana/thema/vmux"
)

//go:embed example.cue cue.mod/module.cue
var LocalSchemaFS embed.FS

func main() {
	ctx := cuecontext.New()
	rt := thema.NewRuntime(ctx)

	// this returns a "not a lineage" error
	_, err := Lineage(rt)
	print(err)

	// This is the manual version of Lineage() above
	bi, err := load.InstanceWithThema(LocalSchemaFS, "")
	exitIf(err)
	v := ctx.BuildInstance(bi)

	print(v)

	// get the lineage
	lineage := v.LookupPath(cue.ParsePath("lineage"))
	print(lineage)
	lin, err := thema.BindLineage(lineage, rt)
	exitIf(err)

	// Getting "not a lineage" even though
	sch, err := lin.Schema(thema.SV(0, 0))
	print(err)

	exampleJSON, _ := ioutil.ReadFile("example.json")
	exdata, err := vmux.NewJSONCodec("example.json").Decode(ctx, exampleJSON)
	print(err)

	_, err = sch.Validate(exdata)
	print(err)
}

func print(thing interface{}) {
	if thing != nil {
		switch thing := thing.(type) {
		case error:
			fmt.Println(thing)
		default:
			fmt.Printf("%#v\n", thing)
		}
	}
}

func exitIf(err error) {
	if err != nil {
		panic(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
