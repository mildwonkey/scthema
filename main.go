package main

import (
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"
)

const dashFile = "dashboard.json"

func main() {
	ctx := cuecontext.New()
	rt := thema.NewRuntime(ctx)
	// Get and decode the dashboard json
	d, _ := os.ReadFile(dashFile)
	dashdata, _ := vmux.NewJSONCodec(dashFile).Decode(ctx, d)

	// Get the lineage, which represents all schemas in the local .cue file
	lin, err := Lineage(rt)
	exitIf(err)

	// Get schema [0,0] and validate
	schema00, _ := lin.Schema(thema.SyntacticVersion{0, 0})
	_, err = schema00.Validate(dashdata)
	if err == nil {
		// This should fail; the json data has a field that was added in 0,1
		fmt.Println("that wasn't supposed to work")
		os.Exit(1)
	} // else task failed succesfully

	// we can explicitly grab schema [0,1], that should be valid
	schema01, err := lin.Schema(thema.SyntacticVersion{0, 1})
	exitIf(err)
	_, err = schema01.Validate(dashdata)
	exitIf(err)

	// but really, we should just let thema tell us what schema fits:
	i := lin.ValidateAny(dashdata)
	if i == nil {
		// if i == nil the bytes do not conform to any schema
		// https://github.com/grafana/thema/issues/156
		fmt.Println("no valid schema found for bits")
	} else {
		fmt.Printf("provided dashboard conforms to schema version %s\n", i.Schema().Version().String())
	}

	// yank the title out of the dashboard json so we can test the lens later.
	i01, _ := i.Translate(thema.SyntacticVersion{0, 1}) // no lacunas; we already know this data fits the schema
	origTitleStr, err := i01.Underlying().LookupPath(cue.ParsePath("title")).String()
	exitIf(err)

	// schema 1,0 adds a "header" field and has a lens to copy the existing "title" to "header"
	start := time.Now()
	// NOTE: i.Translate panics if the target schema does not exist
	// https://github.com/grafana/thema/issues/151
	// i.Translate also panics if lenses are empty or misisng.
	i10, _ := i.Translate(thema.SyntacticVersion{1, 0})
	duration := time.Since(start)
	fmt.Printf("translated to %s in %s seconds\n", i10.Schema().Version().String(), duration)

	// this is all working now!
	headerStr, err := i10.Underlying().LookupPath(cue.ParsePath("header")).String()
	if err != nil {
		fmt.Printf("error looking up header: %s\n", err)
	} else {
		if headerStr != origTitleStr {
			fmt.Printf("title and header should be the same: %q != %q\n", origTitleStr, headerStr)
		}
	}

	// But hey, who cares about instances and bindings! Let's get to the go types!
	var dash Dashboard
	i10.Underlying().Decode(&dash)
	headerStr = dash.Header
	if headerStr != origTitleStr {
		fmt.Printf("title and header should be the same: %q != %q\n", origTitleStr, headerStr)
	}
}

// don't you judge me earl
func print(in interface{}) {
	if in != nil {
		switch in.(type) {
		case error:
			fmt.Println(in)
		default:
			fmt.Printf("%#v\n", in)
		}
	}
}

func exitIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
