package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/thema"
	"github.com/grafana/thema/load"
	"github.com/grafana/thema/vmux"
)

const dashFile = "dashboard.json"

func main() {
	ctx := cuecontext.New()
	rt := thema.NewRuntime(ctx)

	// This loads the core kinds from kindsys plus the dashboard_kind.cue
	bi, err := load.InstanceWithThema(LocalSchemaFS, "")
	exitIfErr(err)
	v := ctx.BuildInstance(bi)

	// get the lineage
	lineage := v.LookupPath(cue.ParsePath("lineage"))
	lin, err := thema.BindLineage(lineage, rt)
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
	} else {
		// task failed succesfully
	}

	// we can explicitly grab schema [0,1], that should be valid
	schema01, _ := lin.Schema(thema.SyntacticVersion{0, 1})
	_, err = schema01.Validate(dashdata)
	exitIfErr(err)

	// but really, we should just let thema tell us what schema fits:
	i := lin.ValidateAny(dashdata)
	if i == nil {
		// this is weird. there's no error message; if i == nil the bytes do not conform to any schema
		fmt.Println("no valid schema found for bits")
	}

	// TODO (to verify below): yank the title out of the dashboard json
	i01, _ := i.Translate(thema.SyntacticVersion{0, 1})
	title := i01.Underlying().Attribute("title")
	print(title.Contents()) // ""

	/*
		Nothing below this comment is quite working the way I expect
		Figure out missing "title" above!
	*/

	// schema v 0,3 renames "title" to "header"; I expected a lacuna beofre the
	// lens was written but there was no output. Is this because i'ts optional?
	// Will try again with schema 0-1
	i02, lacunas := i.Translate(thema.SyntacticVersion{0, 2})
	print(lacunas) // no lacuans in this example, even without the lenses. maybe a major version thing?

	// TODO: there's no lacuna (maybe because it's an optional value), so figure out how
	// to pull the "header" value and confirm that the (now-written) lens is
	// working
	print(i02.Underlying().Attribute("title")) // should not exist in a 0,2 dashboard
	// output: @title()
	// hmmmmmmmmmmmmm

	// the comment on Attribute says any methods on the returned attr should
	// result in an error, but this doesn't return an error - how would I know
	// that "title" doesn't exist in this schema02 instance?
	title = i02.Underlying().Attribute("title")
	print(title.Contents()) // ""

	print(i02.Underlying().Attribute("header")) // should match the earlier "title" field
	// output: @header()
	// hmmmmmmmmmmmmm
	header := i02.Underlying().Attribute("header")
	print(header.Contents()) // ""

	// TODO: GO TYPES
	// why am I hitting Underlying()? Where are my go types?

	// TODO: similar changes, but from schema [0,X] to [1,X]
	// maybe we only get lacunas for required changes?

}

func print(in interface{}) {
	fmt.Printf("%#v\n", in)
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
