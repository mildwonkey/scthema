package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/thema"
	"github.com/grafana/thema/load"
)

func main() {
	ctx := cuecontext.New()
	rt := thema.NewRuntime(ctx)

	// ignoring everything kindsys: just thema & .cue files
	// This is working towards the loading pattern in grafana's plugindef.go
	// however i think is where weird filepath walking / mod prefixing happens in grafana
	bi, err := load.InstanceWithThema(LocalSchemaFS, "dashboard_kind.cue")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v := ctx.BuildInstance(bi)

	// ?????
	lin, err := thema.BindLineage(v, rt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("latest lineage version: %v\n", lin.Latest().Version())
}

func print(in interface{}) {
	fmt.Printf("%#v\n", in)
}

/* Kindsys version
// The package name must be "kindsys" or "" if not passing an overlayFS; the
// difference is unclear. We want the kindsys framework plus our local
// dashboard_kind.cue, so we add the local filesystem with cue files
// embedded. The relpath must match the package name from the cue files; the
// package name must be the same as the relpath or ""

// "relpath" here is the package name in the local .cue files
cv, err := kindsys.BuildInstance(ctx, "scthema", "scthema", LocalSchemaFS)
if err != nil {
	fmt.Println(err.Error())
	os.Exit(1)
}
// print(cv)

// This would be sufficient if we were not loading local CUE files
// fw := kindsys.CUEFramework(ctx)

// still "not a lineage"
lin, err := thema.BindLineage(cv, rt, nil)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
} else {
	print(lin)
}

// STILL still "not a lineage"
raw := rt.Context().BuildInstance(cv.BuildInstance())
lin, err = thema.BindLineage(raw, rt)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
} else {
	print(lin)
}
*/
