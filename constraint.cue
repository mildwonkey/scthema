package scthema

import "github.com/grafana/kindsys"

// NOTE: not sure where "package kind" comes from as described below - maybe it
// should match the package name of this file?

// In each child directory, the set of .cue files with 'package kind'
// must be an instance of kindsys.Core - a declaration of a core kind.
kindsys.Core