package rspns

import "github.com/graaphscom/monogo/compoas"

var ViolationErrorResponseSchema = compoas.Schema{
	Type: "object",
	Properties: map[string]compoas.Schema{
		"Key":     {Type: "string"},
		"Message": {Type: "string"},
		"Violations": {
			Type: "array",
			Items: &compoas.Schema{
				Type: "object",
				Properties: map[string]compoas.Schema{
					"Field":   {Type: "string"},
					"Message": {Type: "string"},
				},
			},
		},
	},
}

var ErrorResponseSchema = compoas.Schema{
	Type: "object",
	Properties: map[string]compoas.Schema{
		"Key":     {Type: "string"},
		"Message": {Type: "string"},
	},
}

func ListResponseSchema(items *compoas.Schema) *compoas.Schema {
	return &compoas.Schema{
		Type: "object",
		Properties: map[string]compoas.Schema{
			"Content":       {Type: "array", Items: items},
			"TotalElements": {Type: "integer"},
			"Size":          {Type: "integer"},
			"TotalPages":    {Type: "integer"},
			"Number":        {Type: "integer"},
		},
	}
}

var ResponseHTTP400 = compoas.Response{
	Description: "Bad request. Returns validation violations",
	Content: map[string]compoas.MediaType{
		"application/json": {Schema: &compoas.Schema{Ref: "#/components/schemas/ViolationError"}},
	},
}

var ResponseHTTP404 = compoas.Response{
	Description: "Not found. Requested resource not found",
	Content: map[string]compoas.MediaType{
		"application/json": {Schema: &compoas.Schema{Ref: "#/components/schemas/Error"}},
	},
}
