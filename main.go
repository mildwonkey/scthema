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

	// TODO: write actual lenses to see how translate works when changes are required
	i, lacunas := i.Translate(thema.SyntacticVersion{0, 2})
	print(lacunas)

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
