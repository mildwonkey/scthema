package main

import (
	"fmt"
	"io/ioutil"
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

	lin, err := Lineage(rt)
	exitIfErr(err)

	// Get and decode the dashboard json
	dash, _ := ioutil.ReadFile(dashFile)
	dashdata, _ := vmux.NewJSONCodec(dashFile).Decode(ctx, dash)

	// Get schema [0,0] and validate (should fail)
	schema00, _ := lin.Schema(thema.SyntacticVersion{0, 0})
	_, err = schema00.Validate(dashdata)
	if err == nil {
		fmt.Println("that wasn't supposed to work")
		os.Exit(1)
	} // else task failed succesfully

	// we can explicitly grab schema [0,1], that should be valid
	schema01, _ := lin.Schema(thema.SyntacticVersion{0, 1})
	_, err = schema01.Validate(dashdata)
	exitIfErr(err)

	// but really, we should just let thema tell us what schema fits:
	i := lin.ValidateAny(dashdata)
	if i == nil {
		// this is unintuitive. there's no error message; if i == nil the bytes
		// do not conform to any schema
		fmt.Println("no valid schema found for bits")
	} else {
		fmt.Printf("provided dashboard conforms to schema version %s\n", i.Schema().Version().String())
	}

	// yank the title out of the dashboard json
	i01, _ := i.Translate(thema.SyntacticVersion{0, 1}) // no lacunas; we already know this data fits the schema
	origTitleStr, err := i01.Underlying().LookupPath(cue.ParsePath("title")).String()
	exitIfErr(err)
	fmt.Printf("original title: %s\n", origTitleStr)

	// schema 0,2 renames "title" to "header"; I expected a lacuna beofre the
	// lens was written but there was no output. Is this because it's optional?
	//
	// This may not be a backwards compatible change by thema's definition, but
	// it's just adding an optional field. I think it's fine.
	fmt.Println("translating to 0,2")
	start := time.Now()
	i02, lacunas := i.Translate(thema.SyntacticVersion{0, 2})
	duration := time.Since(start)
	fmt.Printf("translated to %s in %s seconds\n", i02.Schema().Version().String(), duration)
	print(lacunas.AsList()) // no lacunas in this example, even without the lenses. maybe a major version thing?

	// title is unchanged
	titleStr, err := i02.Underlying().LookupPath(cue.ParsePath("title")).String()
	exitIfErr(err)
	if titleStr != origTitleStr {
		fmt.Println("title should not have changed") // and it won't
		os.Exit(1)
	}

	// new header ~is~ SHOULD BE the same as the original title
	headerStr, err := i02.Underlying().LookupPath(cue.ParsePath("header")).String()
	if err != nil {
		fmt.Println("header not found :(((")
	} else {
		if headerStr != titleStr {
			fmt.Println("title and header should be the same after the translation to 0,2")
		}
	}

	// lazy
	os.Exit(0)

	// TODO: similar changes, but from schema [0,X] to [1,X]
	// maybe we only get lacunas for required changes?
	// Maybe that's why lenses weren't working in the minor changes?

	// Ok, that's it for what I'm now calling a "sequence" translation. Now
	// let's try a "schema" translation - 0,0 to 1,0
	fmt.Println("translating to 1,0")
	start = time.Now()
	// NOTE: i.Translate panics if the target schema does not exist :(
	i10, lacunas := i.Translate(thema.SyntacticVersion{1, 0})
	duration = time.Since(start)
	fmt.Printf("translated to %s in %s seconds\n", i10.Schema().Version().String(), duration)
	print(lacunas.AsList()) // maybe maybe maybe

	headerStr, err = i10.Underlying().LookupPath(cue.ParsePath("header")).String()
	if err != nil {
		fmt.Printf("header string: %s\n", headerStr)
	} else {
		fmt.Println("I did not want an error there")
	}
}

func print(in interface{}) {
	fmt.Printf("%#v\n", in)
}

func exitIfErr(err error) {
	if err != nil {
		panic(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
