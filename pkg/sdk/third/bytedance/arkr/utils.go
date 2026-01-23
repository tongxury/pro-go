package arkr

import "github.com/invopop/jsonschema"

const (
	DefaultModel = Seedance15Pro

	Seedance15Pro = "doubao-seedance-1-5-pro-251215"
	Seed16        = "doubao-seed-1-6-251015"
)

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
