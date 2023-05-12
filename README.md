# what's all this then? 

Objectives:
- [x] use a cue file (dashboard_kind.cue) to vet a json blob's schema (dashboard.json)
- [x] confirm that kindsys is happy with the above (Core kind; ~constraint.cue~)
- [x] Inspect the schema lineages (thema go library)
- [ ] experiment with lineages/schemas/lenses/lacunas 
    - add a lineage, schema, and bidirectional lenses for each

Current status:

* `thema` codegen is working  

Note: we need to identify the path to the lineage within the cue file so thema knows what to work with:
```-p lineage```

There isn't a reserved keyword for the lineage path; it can be called anything (lin, lineage, linmanualmiranda, etc)

Use `LookupPath` to extract the lineage in code:
```lineage := v.LookupPath(cue.ParsePath("lineage"))```


---
_random notes_

straight cue

* constraint.cue: identified as not needed; it has been (or should be) removed from grafana/kindsys
* `cue vet` is happy (dashboard.json && dashboard_kind.cue)
    * TODO: `cue vet dashboard_kind.cue dashboard.json` errors after the thema import was added (for thema.Lineage)

```bash
cue vet dashboard_kind.cue dashboard.json
import failed: cannot find package "github.com/grafana/thema":
    ./dashboard_kind.cue:7:2
```

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
* `load.InstanceWithThema` requires a cue.mod
* constraint.cue causes issues with thema: import failed: cannot find package "github.com/grafana/kindsys"
    * go error or cue error?
* need to include the kindsys.Core

General notes/findings(?)

- `cue.mod` directory is ignored by kindsys loading functions; it's _required_ when using thema but I'm not sure why
- very confusing: cue is declarative, but thema depends on the _order_ of sequences - whre does the schema version go? what's the link between "unnumbered schemas" and "dashboard schema version 36"? shouldn't the schema version be declared as well? 
