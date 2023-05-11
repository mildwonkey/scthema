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
	rt := thema.NewRuntime(ctx)
	lin, err := Lineage(rt)
	if err != nil {
		t.Fatal(err)
	}

	sch, _ := lin.Schema(thema.SV(0, 0))

	data, _ := vmux.NewJSONCodec("dashboard.json").Decode(ctx, dashJSON)
	_, err = sch.Validate(data)
	t.Fatal(err)
}

// https://github.com/grafana/thema/blob/main/docs/go-quickstart.md this gives
// the same "not a lineage" error - does the example work as written? why
// doesn't our stuff ever seem to work with thema?
