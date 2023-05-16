# scthema: experiments with cue, thema and schemas
This repo is an experiment writing schemas and lenses with thema. 

---

## Writing lineages

What's the proper format (cue path?) of a thema lineage? The go types and bindings generate just fine, but I get a "not a lineage" error from the generated Lineage() function. I fixed this by adding a lineage path and extracting the lineage with the path:

```lineage := v.LookupPath(cue.ParsePath("lineage"))```

Note: an invalid module name in module.cue will also result in the same error; is there a better way to capture that?

How do lenses work? What's the format for them? I'm not having luck at the moment. 

## "major" schema version changes

What warrants a new schema (e.g. 1,0 -> 2,0)?
Removing a previously-required field

Schema changes that result in "must be backwards incompatible" errors:
Making an optional field required 
Adding a required field