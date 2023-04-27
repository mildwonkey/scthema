package main

import (
	"embed"
)

//go:embed cue.mod/module.cue *.cue
var LocalSchemaFS embed.FS

//cue.mod is required when using thema's load.InstanceWithThema; this was
//sufficient for kindsys (maybe)
//
//notgo:embed *.cue
