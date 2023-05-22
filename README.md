# what's all this then? 

Objectives:
- [x] use a cue file (dashboard_kind.cue) to vet a json blob's schema (dashboard.json)
- [x] confirm that kindsys is happy with the above (Core kind; ~constraint.cue~)
- [x] inspect the schema lineages (thema go library)
- [ ] experiment with go bindings 
- [x] experiment with go types
    - schemas are complicated, and I'm not in love with pointers for so many types, even though it makes sense. this is definitely weird personal bias
- [ ] experiment with lineages/schemas/lenses/lacunas 
    - add a lineage, schema, and bidirectional lenses for each
    - 0,1 -> 0,2 lens didn't appear to do anything: 
    - renaming 0,2 - 1,0 got me this: schema 1.0 must be backwards incompatible with schema 0.1

Notes:
- the thema cli is not (currently?) the canonical method of generating code for grafana.
    - [TODO] what is the right way and how do external developers use it? 
- Validate() output is frequently misleading. See [grafana/thema #142](https://github.com/grafana/thema/issues/142)

TODOs/side quests - anything mentioned above, plus:
- real dashboard examples as testdata (invalid, valid, valid for different schemas)
- benchmark: real ("real" but manually updated to fit schema 0,0) dashboard schema Validate, Translate funcs
- testdata/www contains random panels from grafana's devenv directory. none are valid. copy into testdata/valid/XXX/* (XXX=min schema version) then make valid for tests.
- What's the correct way to split a single kind.cue into multiple files? The schemas are unweildy
- this, but for every dashboard: https://github.com/grafana/grafana/blob/main/public/app/features/dashboard/state/DashboardMigrator.ts

------
_random notes_

Thema codegen only supports one .cue definition per directory. 

Note: if you prefix the lineage with a key (like the dashboard_kind.cue file in the grafana codebase has at the moment), we need to identify the path to the lineage within the cue file so thema knows what to work with:
```-p lineage```

snippet:
```
lineage: thema.#Lineage  // here "lineage" is the path
lineage: schemas [ ... etc ]
```

There isn't a reserved keyword for the lineage path; it can be called anything (lin, lineage, linmanualmiranda, bo etc)

Use `LookupPath` to extract the lineage in code:
```lineage := v.LookupPath(cue.ParsePath("lineage"))```

* `load.InstanceWithThema` requires a cue.mod
* constraint.cue causes issues with thema: import failed: cannot find package "github.com/grafana/kindsys"
    * go error or cue error?
    * need to include the kindsys.Core
* You will get a "not a lineage" error if the cue.mod/module.cue package name does not match
* You will spend a _lot_ of time stuck on that.
* Just so much time. 

Thema's translate() function is quite slow, 5-15 seconds per call 
```
translating to 0,2
translation took 10.410393416s
```

---

cue specific

* constraint.cue: identified as not needed; it has been (or should be) removed from grafana/kindsys
* `cue vet` is happy (dashboard.json && dashboard_kind.cue)
    * TODO: `cue vet dashboard_kind.cue dashboard.json` errors after the thema import was added (for thema.Lineage)
    - found this note in the thema docs, it should take care of that: https://github.com/grafana/thema/blob/main/docs/go-mapping.md#optional-populate-thema-as-a-cue-dependency

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


General notes/findings(?)

- `cue.mod` directory is ignored by kindsys loading functions; it's _required_ when using thema but I'm not sure why
- cue is declarative, but thema depends on the _order_ of sequences - whre does the schema version go? what's the link between "unnumbered schemas" and "dashboard schema version 36"? shouldn't the schema version be declared as well? (some of this has been improved in thema since writing this)
- `lineage` isn't a keyword to thema; need to pass the actual path. Is `seqs`? Are they all arbitrary keys? That seems fragile; we might need to add stricter checking for that (if there's default or reserved attribute for lineage, we don't need to write `v.LookupPath(cue.ParsePath($LINLINLIN))` every time - maybe that's fine, but it seems like it would be hard to troubleshoot, and hard for third party app developers. Maybe that's something for the SDK)