
what's all this then? 

straight cue

* constraint.cue: ignore for now; must figure out cue loading module ~stuff
* `cue vet` is happy (dashboard.json && dashboard_kind.cue)

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
* load.InstanceWithThema requires a cue.mod
* constraint.cue causes issues with thema: import failed: cannot find package "github.com/grafana/kindsys"
    * go error or cue error?
* still getting `not a lineage (instance root)`
* need to include the kindsys.C

```
Error: not a lineage (instance root): required field is optional in subsumed value: joinSchema (and 1 more errors)
Did you forget to pass a CUE path with -p?
```

TODO
- `constraint.cue` - I have no idea if this is doing anything, or what package names it should use. 
- `cue.mod/` - unused & unclear if it'll be relevant in this case

findings(?)
- kindsys formatted schema isn't thema-able on it's own (using the cli)
- `cue.mod` directory is not necessary in this example
- very confusing: cue is declarative, but thema depends on the order of sequences - whre does the schema version go? what's the link between "unnumbered schemas" and "dashboard schema version 36"? shouldn't the schema version be declared as well? 

