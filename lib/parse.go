package lib

import (
	"encoding/json"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
)

func convert(name string, schema *openapi_v2.Schema, depth int) graphql.Output {
	if depth == 1 {
		return anyType
	}
	schemaType := getType(schema)
	switch schemaType {
	case "anyType":
		return anyType
	case "string":
		return graphql.String
	case "integer":
		return graphql.Int
	case "number":
		return graphql.Float
	case "boolean":
		return graphql.Boolean
	case "array":
		items := schema.GetItems().GetSchema()
		if items != nil {
			return graphql.NewList(createGraphQL(name, items[0], depth-1))
		}
	case "object":
		return createGraphQL(name, schema, depth-1)
	}

	return nil
}

func createGraphQL(name string, schema *openapi_v2.Schema, depth int) *graphql.Object {
	// Check if the type has been visited before, if so, return the cached object.
	if depth <= 0 {
		return graphql.NewObject(graphql.ObjectConfig{
			Name:   getValidGraphQLName(name),
			Fields: graphql.Fields{"anyField": &graphql.Field{Type: anyType}},
		})
	}
	graphqlFields := graphql.Fields{}
	if schema.GetXRef() != "" {
		addRef(schema)
	}
	if schema.Properties != nil {
		for _, property := range schema.Properties.AdditionalProperties {
			log.Info().Msgf("开始处理 %d 层级 %s : %s", depth, name, property.Name)
			if property.Value.GetXRef() != "" {
				addRef(property.Value)
			}
			deal(property.Name, property.Value, graphqlFields, depth)
		}
	}

	if schema.AdditionalProperties != nil {
		shc := schema.AdditionalProperties.GetSchema()
		tmpName := "addName"
		deal(tmpName, shc, graphqlFields, depth)
	}
	if len(graphqlFields) == 0 {
		log.Warn().Msgf("有未处理到的情况,请检查")
	}

	objectType := graphql.NewObject(graphql.ObjectConfig{
		Name:   getValidGraphQLName(name),
		Fields: graphqlFields,
	})
	return objectType
}

func deal(name string, shc *openapi_v2.Schema, graphqlFields graphql.Fields, depth int) {
	fieldType := convert(name, shc, depth)
	currentPropertyName := name
	graphqlFields[name] = &graphql.Field{
		Type: fieldType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Get the source object from the parent.
			s1, ok := p.Source.(gjson.Result)
			var marshal []byte
			if ok {
				marshal = []byte(s1.Raw)
			} else {
				marshal, _ = json.Marshal(p.Source)
			}
			return gjson.GetBytes(marshal, currentPropertyName), nil
		},
	}
}

func addRef(schema *openapi_v2.Schema) {
	definitionKey := schema.GetXRef()[14:]
	shc := findResource(Document, definitionKey)
	if schema.GetProperties().GetAdditionalProperties() == nil {
		add := shc.GetProperties().GetAdditionalProperties()
		prop := openapi_v2.Properties{AdditionalProperties: add}
		schema.Properties = &prop
	} else {
		schema.Properties.AdditionalProperties = append(schema.Properties.AdditionalProperties, shc.GetProperties().GetAdditionalProperties()...)
	}
}
