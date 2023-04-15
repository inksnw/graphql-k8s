package lib

import (
	"encoding/json"
	"fmt"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/tidwall/gjson"
)

func coerceAny(value interface{}) interface{} {
	v, ok := value.(gjson.Result)
	if !ok {
		return value
	}

	switch v.Type {
	case gjson.JSON:
		if v.IsObject() {
			return v.Map()
		}
		if v.IsArray() {
			var array []map[string]any
			err := json.Unmarshal([]byte(v.Raw), &array)
			if err != nil {
				return nil
			}
			return array
		}
	case gjson.String:
		return v.Str
	case gjson.Number:
		return v.Int()
	case gjson.True:
		return v.Bool()
	}
	return value
}

var anyType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "anyType",
	Description: "fake any type",
	Serialize:   coerceAny,
	ParseValue:  coerceAny,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		fmt.Println("处理到这里2")
		return valueAST
	},
})

func getType(schema *openapi_v2.Schema) string {
	schemaTypes := schema.GetType().GetValue()
	if schemaTypes != nil {
		return schemaTypes[0]
	}
	if schema.GetItems() != nil {
		return "array"
	}
	if schema.GetProperties() != nil && len(schema.GetProperties().GetAdditionalProperties()) > 0 {
		return "object"
	}
	if schema.GetAdditionalProperties() != nil && schema.GetAdditionalProperties().GetSchema() != nil {
		return "object"
	}
	return "anyType"
}
