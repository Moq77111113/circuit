package form

import (
	"reflect"

	"github.com/moq77111113/circuit/internal/schema"
)
// findNodeAndValue locates the Node and its corresponding reflect.Value based on the provided path.
func findNodeAndValue(nodes []schema.Node, path schema.Path, rootValue reflect.Value) (*schema.Node, reflect.Value) {
	fieldPath := path.String()
	if fieldPath == "" {
		return nil, reflect.Value{}
	}

	for i := range nodes {
		if nodes[i].Name == fieldPath {
			fieldValue := rootValue.FieldByName(nodes[i].Name)
			return &nodes[i], fieldValue
		}
	}

	segments := extractPathSegments(fieldPath)
	if len(segments) == 0 {
		return nil, reflect.Value{}
	}
	currentNodes := nodes
	currentValue := rootValue

	for i, seg := range segments {
		var foundNode *schema.Node
		for j := range currentNodes {
			if currentNodes[j].Name == seg.name {
				foundNode = &currentNodes[j]
				break
			}
		}

		if foundNode == nil {
			return nil, reflect.Value{}
		}

		fieldValue := currentValue.FieldByName(foundNode.Name)
		if !fieldValue.IsValid() {
			return nil, reflect.Value{}
		}

		if seg.hasIndex && foundNode.Kind == schema.KindSlice {
			if seg.index < 0 || seg.index >= fieldValue.Len() {
				return nil, reflect.Value{}
			}
			fieldValue = fieldValue.Index(seg.index)
		}

		if i == len(segments)-1 {
			return foundNode, fieldValue
		}

		if foundNode.Kind == schema.KindStruct || foundNode.Kind == schema.KindSlice {
			if len(foundNode.Children) > 0 {
				currentNodes = foundNode.Children
				currentValue = fieldValue
			} else {
				return nil, reflect.Value{}
			}
		} else {
			return nil, reflect.Value{}
		}
	}

	return nil, reflect.Value{}
}
