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
			},
		},
		{
			version: [0, 1]
			schema:
			{
				title: string
				header?: string // new, optional field ->  no error
			},
		},
		// {
		// 	version: [1, 0]
		// 	schema:
		// 	{
		// 		// title: string <-- this is the only *backwards* incompatible change
		// 		header: string // optional in 0,1; required in schema 1,0
		// 		newRequiredField: string
		// 	},
		// },
	]
	// lenses: [
	// 	{
	// 		to: [0, 1]
	// 		from: [0, 0]
	// 		input: _
	// 		result: {
	// 			header: input.title
	// 		}
	// 	},
	// 	{
	// 		to: [1, 0]
	// 		from: [0, 1]
	// 		input: _
	// 		result: {
	// 			newRequiredField: "default value" // I'm not sure if this is valid, worth a go.
	// 		}
	// 	},
	// ]
}
