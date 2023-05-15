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
	} else {
		fmt.Printf("found matching version %s\n", i.Schema().Version().String())
	}

	// TODO (to verify below): yank the title out of the dashboard json
	i01, _ := i.Translate(thema.SyntacticVersion{0, 1}) // no lacunas; we already know this data fits the schema
	origTitleStr, err := i01.Underlying().LookupPath(cue.ParsePath("title")).String()
	exitIfErr(err)
	fmt.Printf("original title: %s\n", origTitleStr)

	// schema 0,2 renames "title" to "header"; I expected a lacuna beofre the
	// lens was written but there was no output. Is this because i'ts optional?
	// Will try again with schema 0-1
	fmt.Println("translating to 0,2")
	i02, lacunas := i.Translate(thema.SyntacticVersion{0, 2})
	fmt.Printf("translated to %s", i02.Schema().Version().String())
	print(lacunas.AsList()) // no lacunas in this example, even without the lenses. maybe a major version thing?

	// the comment on Attribute says any methods on the returned attr should
	// result in an error, but this doesn't return an error - how would I know
	// that "title" doesn't exist in this schema02 instance?
	titleStr, err := i02.Underlying().LookupPath(cue.ParsePath("title")).String() // should not exist in a 0,2 dashboard
	if err != nil {
		fmt.Printf("I thought that should have errored, but instead we got %#v\n", titleStr)
		// try dehydrating the instance
		str, err := i02.Dehydrate().Underlying().LookupPath(cue.ParsePath("title")).String()
		if err != nil {
			print(err)
		} else {
			fmt.Printf("dehydrated title: %s\n", str)
		}
	}

	headerStr, err := i02.Underlying().LookupPath(cue.ParsePath("header")).String()
	exitIfErr(err)
	if headerStr != origTitleStr {
		fmt.Println("title and header should be the same")
		os.Exit(1)
	}

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
