package openapi

import "github.com/go-openapi/spec"

var Spec = &spec.Swagger{
	SwaggerProps: spec.SwaggerProps{
		Swagger: "2.0",
		Consumes: []string{
			"application/json",
			"application/xml",
		},
		Paths: &spec.Paths{
			Paths: map[string]spec.PathItem{},
		},
		Definitions: spec.Definitions{},
	},
}
