package walk

import (
	"errors"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

var ErrSkipChildren = errors.New("skip children")

// WalkConfig holds configuration for tree walking.
type WalkConfig struct {
	MaxDepth int
	BasePath path.Path
}

// WalkOption configures a Walker.
type WalkOption func(*WalkConfig)

// WithMaxDepth sets the maximum depth for traversal.
func WithMaxDepth(depth int) WalkOption {
	return func(cfg *WalkConfig) {
		cfg.MaxDepth = depth
	}
}

// WithBasePath sets the base path to prepend to all paths during traversal.
func WithBasePath(p path.Path) WalkOption {
	return func(cfg *WalkConfig) {
		cfg.BasePath = p
	}
}

// Walker traverses an AST and calls visitor methods.
type Walker struct {
	visitor Visitor
	config  WalkConfig
}

// NewWalker creates a new Walker with the given visitor and options.
func NewWalker(v Visitor, opts ...WalkOption) *Walker {
	w := &Walker{
		visitor: v,
		config:  WalkConfig{MaxDepth: -1},
	}
	for _, opt := range opts {
		opt(&w.config)
	}
	return w
}

// Walk traverses the tree and calls visitor methods.
func (w *Walker) Walk(tree *node.Tree, state any) error {
	ctx := NewContext(tree, state, w.config.BasePath)
	return w.walkNodes(tree.Nodes, ctx)
}

// WalkWithContext traverses the tree with an additional immutable context.
func (w *Walker) WalkWithContext(tree *node.Tree, state any, context any) error {
	ctx := NewContext(tree, state, w.config.BasePath)
	ctx.Context = context
	return w.walkNodes(tree.Nodes, ctx)
}

func (w *Walker) walkNodes(nodes []node.Node, ctx *VisitContext) error {
	for i := range nodes {
		nodeCtx := *ctx
		nodeCtx.Path = ctx.Path.Child(nodes[i].Name)

		if err := w.walkNode(&nodes[i], &nodeCtx); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkNode(n *node.Node, ctx *VisitContext) error {
	if w.config.MaxDepth >= 0 && ctx.Depth > w.config.MaxDepth {
		return nil
	}

	err := w.visit(n, ctx)
	if err != nil {
		if errors.Is(err, ErrSkipChildren) {
			return nil
		}
		return err
	}

	// Recurse into children for structs
	if n.Kind == node.KindStruct && len(n.Children) > 0 {
		childCtx := *ctx
		childCtx.Depth++
		childCtx.Parent = n
		return w.walkNodes(n.Children, &childCtx)
	}

	return nil
}

func (w *Walker) visit(n *node.Node, ctx *VisitContext) error {
	switch n.Kind {
	case node.KindPrimitive:
		return w.visitor.VisitPrimitive(ctx, n)
	case node.KindStruct:
		return w.visitor.VisitStruct(ctx, n)
	case node.KindSlice:
		return w.visitor.VisitSlice(ctx, n)
	}
	return nil
}
