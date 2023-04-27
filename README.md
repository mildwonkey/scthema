# what's all this then? 

Objectives: 
* use a cue file (dashboard_kind.cue) to vet a json blob's schema (dashboard.json)
* confirm that kindsys is happy with the above (Core kind; constraint.cue)
* Inspect the schema lineages (thema go library)

Current status:

thema cli commands work; I needed to modify the cue files (there's a PR in grafana to update them, so I will grab the current cue definition when it's merged)

> Note: we need to include `-p lineage` to the thema cli commands to get past the "not a lineage" errors

The thema library is giving the same "not a lineage" error which was resolved in the CLI with `-p` flag. Next step is figuring out how to do that in code.

---
_random notes_

straight cue

* constraint.cue: ignore for now; must figure out cue loading module ~stuff
* `cue vet` is happy (dashboard.json && dashboard_kind.cue)
    * NOTE: `cue vet dashboard_kind.cue dashboard.json` errors after the thema import was added (for thema.Lineage)

Loading CUE files from multiple sources:
kindsys.BuildInstance  will pull all the kindsys cue files in regardless of how it's called.  
To load local CUE files alongside the kindsys framework:
- pick a package name for the CUE files and stick with it (it does not need to match the go package)
- embed the local CUE files with embed.FS 
- call kindsys.BuildInstance with args:
    * ctx (`*cue.Context`)
    * "package name from local .cue files" 
    * "package name or blank, nothing else works" 
        - I don't know what the difference between those options would be
    * LocalEmbeddedFS (most basic usage: `*.cue`) 
        - module.cue is a place for vendored 
        - not sure how we end up with files in the cue.mod directory, or if it's relevant to this use case
        - perhaps it'd be an alternative to kindsys.BuildInstance - just copy all the kinds into cue.mod/

thema 
* cli not working with kindsyssy files: `thema lineage gen gotypes -l dashboard_kind.cue`
    * note: worked with updated files + `-p` flag
* `load.InstanceWithThema` requires a cue.mod
* constraint.cue causes issues with thema: import failed: cannot find package "github.com/grafana/kindsys"
    * go error or cue error?
* still getting `not a lineage (instance root)`
* need to include the kindsys.Core

```
Error: not a lineage (instance root): required field is optional in subsumed value: joinSchema (and 1 more errors)
Did you forget to pass a CUE path with -p?
```

TODO
- `constraint.cue` - actually figure out what it's doing and how it works; how does it interact with `kindsys`? 

findings(?)

- `cue.mod` directory is ignored by kindsys loading functions; it's _required_ when using thema but I'm not sure why
- very confusing: cue is declarative, but thema depends on the _order_ of sequences - whre does the schema version go? what's the link between "unnumbered schemas" and "dashboard schema version 36"? shouldn't the schema version be declared as well? 

