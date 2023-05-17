package main

import (
	"embed"
)

//go:embed cue.mod/module.cue *.cue
var LocalSchemaFS embed.FS
