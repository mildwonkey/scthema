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
	// bits and bobs to get started
	ctx := cuecontext.New()
	rt := thema.NewRuntime(ctx)
	exampleJSON, _ := ioutil.ReadFile("example.json")
	exdata, _ := vmux.NewJSONCodec("example.json").Decode(ctx, exampleJSON)

	// Using the generated go binding this returns a "not a lineage" error Note:
	// I've leard that this error can happen if module.cue has the wrong package
	// name, though that's not the case here - perhaps thema expects a specific
	// cue path for the Lineage?
	_, err := Lineage(rt)
	print(err)

	// This is the manual version of Lineage() above; the only difference is
	// adding the lookup path
	bi, err := load.InstanceWithThema(LocalSchemaFS, "")
	exitIf(err)
	v := ctx.BuildInstance(bi)

	// get the lineage and a schema, and validate the example.json data
	lineage := v.LookupPath(cue.ParsePath("lineage"))
	lin, err := thema.BindLineage(lineage, rt)
	exitIf(err)

	sch00, err := lin.Schema(thema.SV(0, 0)) // we wouldn't normally hardcode this; just for the example
	exitIf(err)
	i00, err := sch00.Validate(exdata)
	exitIf(err)

	// let's save the title in a var for later comparison. Title is an optional field with no default
	origTitleStr, err := i00.Underlying().LookupPath(cue.ParsePath("title")).String()
	exitIf(err)
	fmt.Printf("Original title: %s\n", origTitleStr)

	// translate the original example.json into the new schema 0,1. 0,1 adds an
	// optional header field; I've attempted to map the title field to the
	// header using a lens.
	i10, _ := i00.Translate(thema.SV(0, 1)) // no lacunas
	_, err = i10.Underlying().LookupPath(cue.ParsePath("title")).String()
	print(err) // #Translate.out.result.result: field not found: title

	// NOTE: still working out what makes a change backwards incompatible by thema's heuristics.
	// Making an optional field required resulted in this error:
	//    schema 1.0 must be backwards incompatible with schema 0.1
	// I had hoped to use a lens to map
	// The above code worked before adding the 1,0 schema, but now I'm
	// getting an error:
	//    schema 1.0 must be backwards incompatible with schema 0.1
	// It works if I remove the `title` field entirely, but making an optional
	// field (with no default) is not backwards incompatible by thema's
	// heuristics at time of writing.

	// Let's try something similar, but make the change in a new schema
	i10, _ = i00.Translate(thema.SV(1, 0))
	_, err = i10.Underlying().LookupPath(cue.ParsePath("title")).String()
	print(err) // #Translate.out.result.result: field not found: title

	titleStr10, err := i10.Underlying().LookupPath(cue.ParsePath("title")).String()
	if err != nil {
		print(err)
	} else if titleStr10 != origTitleStr {
		fmt.Println("title changed")
	}

	reqdStr10, err := i10.Underlying().LookupPath(cue.ParsePath("newRequiredField")).String()
	if err != nil {
		print(err) // field not found: newRequiredField
	} else {
		print(reqdStr10)
	}

	headerStr10, err := i10.Underlying().LookupPath(cue.ParsePath("header")).String()
	if err != nil {
		print(err) // field not found: header
	} else {
		print(headerStr10)
	}

}

// pls don't look below this line ;)
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
		fmt.Println(err)
		os.Exit(1)
	}
}
