package lib

import (
	"encoding/json"
	"fmt"
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
	"github.com/phuslu/log"
	"github.com/tidwall/gjson"
	"strings"
)

func createGraphQL(name string, schema *openapi_v2.Schema, depth int) *graphql.Object {
	validName := getValidName(name)
	if depth <= 0 {
		result := graphql.NewObject(graphql.ObjectConfig{
			Name:   validName,
			Fields: graphql.Fields{"anyField": &graphql.Field{Type: anyType}},
		})
		return result
	}
	graphqlFields := graphql.Fields{}
	if schema.GetXRef() != "" {
		addRef(schema)
	}
	if schema.GetProperties() != nil {
		for _, property := range schema.Properties.AdditionalProperties {
			//log.Debug().Msgf("开始处理 %d 层级 %s : %s", depth, name, property.Name)
			if property.Value.GetXRef() != "" {
				addRef(property.Value)
			}
			fullName := fmt.Sprintf("%s_%s", name, property.Name)
			deal(fullName, property.Value, graphqlFields, depth)
		}
	}

	if schema.GetAdditionalProperties() != nil {
		shc := schema.AdditionalProperties.GetSchema()
		//log.Debug().Msgf("处理AdditionalProperties %d 层级 %s : %s", depth, name, name)
		deal(name, shc, graphqlFields, depth)
	}
	if len(graphqlFields) == 0 {
		log.Warn().Msgf("有未处理到的情况,请检查 %s", name)
		result := graphql.NewObject(graphql.ObjectConfig{
			Name:   validName,
			Fields: graphql.Fields{"anyField": &graphql.Field{Type: anyType}},
		})
		return result
	}

	objectType := graphql.NewObject(graphql.ObjectConfig{
		Name:   validName,
		Fields: graphqlFields,
	})
	return objectType
}

func deal(name string, shc *openapi_v2.Schema, graphqlFields graphql.Fields, depth int) {
	var fieldType graphql.Output

	schemaType := getType(shc)
	switch schemaType {
	case "anyType":
		fieldType = anyType
	case "string":
		fieldType = graphql.String
	case "integer":
		fieldType = graphql.Int
	case "number":
		fieldType = graphql.Float
	case "boolean":
		fieldType = graphql.Boolean
	case "array":
		if depth == 1 {
			//解决graphql强类型,对于有下级结构的必须显示指定字段
			fieldType = anyType
		} else {
			items := shc.GetItems().GetSchema()
			oneType := createGraphQL(name, items[0], depth-1)
			fieldType = graphql.NewList(oneType)
		}

	case "object":
		if depth == 1 {
			//解决graphql强类型,对于有下级结构的必须显示指定字段
			fieldType = anyType
		} else {
			fieldType = createGraphQL(name, shc, depth-1)
		}

	}

	graphqlFields[name] = &graphql.Field{
		Type:    fieldType,
		Resolve: resolve(name),
	}
}

func resolve(name string) graphql.FieldResolveFn {
	currentPropertyName := name

	//解决graphql一个schema下不支持同名的字段
	if strings.Contains(currentPropertyName, "_") {
		list := strings.Split(currentPropertyName, "_")
		currentPropertyName = list[len(list)-1]
	}

	return func(p graphql.ResolveParams) (interface{}, error) {
		// Get the source object from the parent.
		s1, ok := p.Source.(gjson.Result)
		var marshal []byte
		if ok {
			marshal = []byte(s1.Raw)
		} else {
			marshal, _ = json.Marshal(p.Source)
		}
		result := gjson.GetBytes(marshal, currentPropertyName)

		if result.IsArray() {
			var array []map[string]any
			err := json.Unmarshal([]byte(result.Raw), &array)
			return array, err
		}
		return result, nil
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
