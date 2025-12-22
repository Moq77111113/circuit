package breadcrumb

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func RenderBreadcrumb(currentPath path.Path, nodes []ast.Node) g.Node {
	segments := currentPath.Segments()

	var items []g.Node
	items = append(items, renderRootLink())

	if len(segments) == 0 {
		return h.Nav(h.Class("breadcrumb"), g.Group(items))
	}

	for i := range segments {
		seg := segments[i]
		items = append(items, renderSeparator())

		if isIndex(seg) {
			items[len(items)-1] = renderIndexSegment(seg)
			continue
		}

		partialPath := buildPartialPath(segments[:i+1])
		label := humanizeLabel(seg, nodes)
		items = append(items, renderLink(label, partialPath))
	}

	return h.Nav(h.Class("breadcrumb"), g.Group(items))
}

func renderRootLink() g.Node {
	return h.A(
		h.Href("?focus="),
		h.Class("breadcrumb__link breadcrumb__link--root"),
		g.Text("Config"),
	)
}

func renderLink(label string, pathStr string) g.Node {
	return h.A(
		h.Href("?focus="+pathStr),
		h.Class("breadcrumb__link"),
		g.Text(label),
	)
}

func renderSeparator() g.Node {
	return h.Span(
		h.Class("breadcrumb__separator"),
		g.Text(" > "),
	)
}

func renderIndexSegment(idx string) g.Node {
	return h.Span(
		h.Class("breadcrumb__index"),
		g.Text("["+idx+"]"),
	)
}

func buildPartialPath(segments []string) string {
	result := ""
	for i, seg := range segments {
		if i > 0 {
			result += "."
		}
		result += seg
	}
	return result
}

func humanizeLabel(segment string, nodes []ast.Node) string {
	return segment
}

func isIndex(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
