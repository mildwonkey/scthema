package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"testing"

	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"
)

//go:embed dashboard.json
var dashJSON []byte

//go:embed testdata/valid/*.json
var validTestData embed.FS

//go:embed testdata/invalid/*.json
var invalidTestData embed.FS

var (
	data, _  = vmux.NewJSONCodec("dashboard.json").Decode(ctx, dashJSON)
	lin, _   = Lineage(thema.NewRuntime(ctx))
	sch01, _ = lin.Schema(thema.SV(0, 1))
)

func TestValidateDash(t *testing.T) {
	t.Run("valid files", func(t *testing.T) {
		err := fs.WalkDir(validTestData, ".", func(path string, d fs.DirEntry, err error) error {
			// walk the testData filesystem and validate each .json file
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil // skip directories
			}

			f, err := validTestData.ReadFile(path)
			if err != nil {
				return err
			}

			data, err := vmux.NewJSONCodec(d.Name()).Decode(ctx, f)
			if err != nil {
				return err
			}

			_, err = sch01.Validate(data)
			return err
		})
		if err != nil {
			t.Fatal("validate errors") // the errors are verbose
		}
	})

	// this is the least useful test, on it's own: any errors would make it
	// pass, not just validate errors. good enough for this li'l repo, though.
	t.Run("invalid files", func(t *testing.T) {
		err := fs.WalkDir(invalidTestData, ".", func(path string, d fs.DirEntry, err error) error {
			// walk the testData filesystem and validate each .json file
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil // skip directories
			}

			f, err := invalidTestData.ReadFile(path)
			if err != nil {
				return err
			}

			data, err := vmux.NewJSONCodec(d.Name()).Decode(ctx, f)
			if err != nil {
				return err
			}

			_, err = sch01.Validate(data)
			return err
		})
		if err == nil {
			t.Fatal("expected validate errors")
		}
	})
}

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sch01.Validate(data)
	}
}

func BenchmarkValidate_invalid(b *testing.B) {
	sch, _ := lin.Schema(thema.SV(0, 0))
	for i := 0; i < b.N; i++ {
		_, err := sch.Validate(data)
		if err == nil {
			b.Fatal("expected validate errors")
		}
	}
}

func BenchmarkValidateAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lin.ValidateAny(data)
	}
}

func TestTranslate(t *testing.T) {
	inst, _ := sch01.Validate(data)
	// There's no error return from Translate(yet), so we'll just make sure it
	// doesn't panic.
	_, _ = inst.Translate(thema.SV(1, 0))

}

// This only runs once as-is; add `-benchtime=20x` to the command line. (100x
// timed out)
func BenchmarkTranslate(b *testing.B) {
	inst, _ := sch01.Validate(data)
	for i := 0; i < b.N; i++ {
		inst.Translate(thema.SV(1, 0))
	}
}
