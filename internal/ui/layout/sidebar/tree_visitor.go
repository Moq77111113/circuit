package sidebar

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

type TreeVisitor struct {
	currentFocus path.Path
	values       map[string]any
}

func (v *TreeVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*TreeState)
	isActive := ctx.Path.String() == v.currentFocus.String()

	state.Append(renderTreeLeaf(node.Name, ctx.Path, isActive))
	return nil
}

func (v *TreeVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*TreeState)
	isActive := ctx.Path.String() == v.currentFocus.String()
	isExpanded := v.currentFocus.HasPrefix(ctx.Path) || ctx.Depth == 0

	state.Append(renderTreeNode(node.Name, ctx.Path, isActive, isExpanded, true))
	return nil
}

func (v *TreeVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	state := ctx.State.(*TreeState)
	isExpanded := v.currentFocus.HasPrefix(ctx.Path) || ctx.Depth == 0

	state.Append(renderTreeNode(node.Name, ctx.Path, false, isExpanded, true))
	return nil
}

func renderTreeLeaf(name string, nodePath path.Path, isActive bool) g.Node {
	class := "tree-node tree-node--leaf"
	if isActive {
		class += " tree-node--active"
	}

	return h.Div(
		h.Class(class),
		h.A(
			h.Href("?focus="+nodePath.String()),
			h.Class("tree-node__link"),
			g.Text(name),
		),
	)
}

func renderTreeNode(name string, nodePath path.Path, isActive, isExpanded, hasChildren bool) g.Node {
	class := "tree-node"
	if isActive {
		class += " tree-node--active"
	}
	if isExpanded {
		class += " tree-node--expanded"
	}

	var icon string
	if hasChildren {
		if isExpanded {
			icon = "▼ "
		} else {
			icon = "▶ "
		}
	} else {
		icon = "• "
	}

	return h.Div(
		h.Class(class),
		h.A(
			h.Href("?focus="+nodePath.String()),
			h.Class("tree-node__link"),
			g.Text(icon+name),
		),
	)
}
