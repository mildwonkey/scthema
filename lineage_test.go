package main

import (
	_ "embed"
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"
)

//go:embed dashboard.json
var dashJSON []byte

func TestValidateDash(t *testing.T) {
	ctx := cuecontext.New()
	lin, err := Lineage(thema.NewRuntime(ctx))
	if err != nil {
		t.Fatal(err)
	}

	sch, _ := lin.Schema(thema.SV(0, 0))

	data, _ := vmux.NewJSONCodec("dashboard.json").Decode(ctx, dashJSON)
	_, err = sch.Validate(data)
	t.Fatal(err)
}
