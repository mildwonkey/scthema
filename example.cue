package main

import (
	"github.com/grafana/thema"
)

lineage: thema.#Lineage
lineage: {
	name: "example"
	schemas: [
		{
			version: [0, 0]
			schema:
			{
				title: string
				optionalTitle?: string
			},
		},
		{
			version: [0, 1]
			schema:
			{
				title: string
				header?: string // new, optional field ->  no error
				
				// Can't make a previously optional field required in the same major version
				//     required field is optional in subsumed value: optionalTitle
				// optionalTitle: string
				optionalTitle?: string
			},
		},
		{
			version: [1, 0]
			schema:
			{
				// title: string --> this schema was not considered a breaking change until i removed this
				header: string // optional in 0,1; required in schema 1,0
				optionalTitle: string // previously optional, now required
				newRequiredField: string
			},
		},
	]
	lenses: [
		{
			to: [0, 1]
			from: [0, 0]
			input: _
			result: {
				header: input.title
			}
		},
		{
			to: [1, 0]
			from: [0, 1]
			input: _
			result: {
				header: input.title
				newRequiredField: "default value"
			}
		},
	]
}
