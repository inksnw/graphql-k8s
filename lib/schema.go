package lib

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"github.com/graphql-go/graphql"
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

type Action struct {
	Kind string
	Shc  *graphql.Schema
}

func k8sResolve(gvr schema.GroupVersionResource) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		var result []map[string]any
		name := p.Args["name"].(string)
		namespace := p.Args["namespace"].(string)
		if namespace == "" {
			namespace = metav1.NamespaceDefault
		}

		if name != "" {
			rv, err := DynamicClient.Resource(gvr).Namespace(namespace).Get(p.Context, name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			result = append(result, rv.Object)
			return result, err
		}
		label := p.Args["label"].(string)
		var option metav1.ListOptions
		if label != "" {
			option = metav1.ListOptions{LabelSelector: label}
		}
		unstructuredList, err := DynamicClient.Resource(gvr).List(p.Context, option)
		for _, i := range unstructuredList.Items {
			result = append(result, i.Object)
		}
		return result, err
	}

}

func GenerateGraphQLSchema(resources []ResourceType, depth int) (ActionMap map[string]*graphql.Schema, err error) {
	ActionMap = make(map[string]*graphql.Schema)
	for _, r := range resources {
		definition := findResource(Document, r.ResourceName)
		Type := createGraphQL(r.Kind, definition, depth)
		gvr := schema.GroupVersionResource{
			Group:    r.Group,
			Version:  r.Version,
			Resource: r.Resource,
		}

		shc, err := graphql.NewSchema(graphql.SchemaConfig{
			Query: graphql.NewObject(graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					r.Kind: &graphql.Field{
						Type:    graphql.NewList(Type),
						Args:    buildQueryArgs(),
						Resolve: k8sResolve(gvr),
					},
				},
			}),
		})
		if err != nil {
			return nil, err
		}

		ActionMap[strings.ToLower(r.Kind)] = &shc
	}

	return ActionMap, err
}

func buildQueryArgs() graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			DefaultValue: "",
			Description:  "The metadata.name of the Pod",
			Type:         graphql.String,
		},
		"namespace": &graphql.ArgumentConfig{
			DefaultValue: "",
			Description:  "The metadata.namespace of the Pod",
			Type:         graphql.String,
		},
		"label": &graphql.ArgumentConfig{
			DefaultValue: "",
			Description:  "The metadata.labels of the Pod",
			Type:         graphql.String,
		},
	}
}
