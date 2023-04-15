package lib

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"unicode"
)

func getValidName(name string) string {
	// Replace any character that is not a letter, digit, or underscore with an underscore.
	// Ensure the first character is an uppercase letter.
	var newName strings.Builder
	newName.Grow(len(name) + 1)
	newName.WriteByte(byte(unicode.ToUpper(rune(name[0]))))

	for i := 1; i < len(name); i++ {
		ch := rune(name[i])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			newName.WriteByte(byte(ch))
		} else {
			newName.WriteByte('_')
		}
	}

	return newName.String()
}

func findResource(document *openapi_v2.Document, resourceName string) *openapi_v2.Schema {
	definitions := document.GetDefinitions()
	for _, namedSchema := range definitions.AdditionalProperties {
		if namedSchema.Name == resourceName {
			return namedSchema.Value
		}
	}
	return nil
}

func GenerateGraphQLSchema(resources []ResourceType, depth int) (*graphql.Schema, error) {

	fields := graphql.Fields{}

	for _, r := range resources {
		definition := findResource(Document, r.ResourceName)
		Type := createGraphQL(r.Kind, definition, depth)
		gvr := schema.GroupVersionResource{
			Group:    r.Group,
			Version:  r.Version,
			Resource: r.Resource,
		}
		fields[r.Resource] = &graphql.Field{
			Type: graphql.NewList(Type),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				unstructuredList, err := DynamicClient.Resource(gvr).List(p.Context, metav1.ListOptions{})
				var result []map[string]any
				for _, i := range unstructuredList.Items {
					result = append(result, i.Object)
				}
				return result, err
			},
		}
	}

	// Create the schema from the query type
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
		}),
	})
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	return &schema, err
}
