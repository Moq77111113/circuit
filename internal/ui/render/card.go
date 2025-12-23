package render

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ui/styles"
)

func RenderStructCard(node ast.Node, nodePath path.Path, values map[string]any) g.Node {
	focusURL := "?focus=" + nodePath.String()
	preview := generatePreview(node, nodePath, values, 3)

	return h.A(
		h.Href(focusURL),
		h.Class(styles.StructCard),
		h.Div(
			h.Class(styles.StructCardHeader),
			h.Span(h.Class(styles.StructCardName), g.Text(node.Name)),
			h.Span(h.Class(styles.StructCardArrow+" "+styles.IconArrowRight)),
		),
		h.Div(
			h.Class(styles.StructCardPreview),
			g.Text(preview),
		),
	)
}

func generatePreview(node ast.Node, nodePath path.Path, values map[string]any, maxFields int) string {
	if len(node.Children) == 0 {
		return ""
	}

	var previews []string
	count := 0

	for _, child := range node.Children {
		if child.Kind != ast.KindPrimitive {
			continue
		}
		if count >= maxFields {
			previews = append(previews, "...")
			break
		}

		childPath := nodePath.Child(child.Name)
		value := values[childPath.String()]

		if value != nil {
			preview := formatPreviewValue(child.Name, value)
			previews = append(previews, preview)
			count++
		}
	}

	if len(previews) == 0 {
		return "No values set"
	}

	result := ""
	for i, p := range previews {
		if i > 0 {
			result += " | "
		}
		result += p
	}
	return result
}

func formatPreviewValue(name string, value any) string {
	switch v := value.(type) {
	case string:
		if len(v) > 20 {
			return fmt.Sprintf("%s: %s...", name, v[:20])
		}
		return fmt.Sprintf("%s: %s", name, v)
	case bool:
		return fmt.Sprintf("%s: %t", name, v)
	case int, int64:
		return fmt.Sprintf("%s: %v", name, v)
	case float64:
		return fmt.Sprintf("%s: %.2f", name, v)
	default:
		return fmt.Sprintf("%s: %v", name, v)
	}
}
