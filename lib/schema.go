package lib

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
	"github.com/phuslu/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"unicode"
)

func getValidGraphQLName(name string) string {
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

func GenerateGraphQLSchema(depth int) (*graphql.Schema, error) {

	Resources := []string{"io.k8s.api.core.v1.Pod"}
	fields := graphql.Fields{}

	for _, r := range Resources {
		definition := findResource(Document, r)
		Type := createGraphQL("Pod", definition, depth)
		fields["pods"] = &graphql.Field{
			Type: graphql.NewList(Type),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pods, err := Clientset.CoreV1().Pods("").List(p.Context, metav1.ListOptions{})
				if err != nil {
					return nil, err
				}
				for idx := range pods.Items {
					pods.Items[idx].Kind = "Pod"
					pods.Items[idx].APIVersion = "v1"
				}
				return pods.Items, err
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
