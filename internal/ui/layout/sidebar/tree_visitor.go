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
	values       ast.ValuesByPath
}

func (v *TreeVisitor) VisitPrimitive(ctx *walk.VisitContext, node *ast.Node) error {
	if ctx.Depth > 0 {
		return nil
	}

	state := ctx.State.(*TreeState)
	isActive := ctx.Path.String() == v.currentFocus.String()

	state.Append(renderTreeLeaf(node.Name, ctx.Path, isActive))
	return nil
}

func (v *TreeVisitor) VisitStruct(ctx *walk.VisitContext, node *ast.Node) error {
	// Only show root-level structs
	if ctx.Depth > 0 {
		return walk.ErrSkipChildren
	}

	state := ctx.State.(*TreeState)
	isActive := ctx.Path.String() == v.currentFocus.String()

	// Root structs are simple links (no expansion/chevron)
	state.Append(renderTreeLeaf(node.Name, ctx.Path, isActive))
	return walk.ErrSkipChildren
}

func (v *TreeVisitor) VisitSlice(ctx *walk.VisitContext, node *ast.Node) error {
	if ctx.Depth > 0 {
		return nil
	}

	state := ctx.State.(*TreeState)
	isActive := ctx.Path.String() == v.currentFocus.String()

	state.Append(renderTreeLeaf(node.Name, ctx.Path, isActive))
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
