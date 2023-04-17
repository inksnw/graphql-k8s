package lib

import (
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	"strings"
)

func AddPrefix(node ast.Node, parentKeys []string) {

	switch n := node.(type) {
	case *ast.Document:
		for _, def := range n.Definitions {
			AddPrefix(def, parentKeys)
		}
	case *ast.OperationDefinition:
		AddPrefix(n.SelectionSet, parentKeys)
	case *ast.SelectionSet:
		for _, sel := range n.Selections {
			selNode := sel.(ast.Node)
			AddPrefix(selNode, parentKeys)
		}
	case *ast.Field:
		name := fmt.Sprintf("%s", n.Name.Value)
		if parentKeys != nil {
			newName := strings.Join(append(parentKeys, n.Name.Value), "_")
			n.Name.Value = newName
		}
		parentKeys = append(parentKeys, name)

		if n.SelectionSet != nil {
			AddPrefix(n.SelectionSet, parentKeys)
		}
	default:
		// do nothing
	}
}
