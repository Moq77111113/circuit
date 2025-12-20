package ast

import (
	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/ast/walk"
)

type (
	Node       = node.Node
	UIMetadata = node.UIMetadata
	Schema     = node.Schema
	Tree       = node.Tree
)

const (
	KindPrimitive = node.KindPrimitive
	KindStruct    = node.KindStruct
	KindSlice     = node.KindSlice

	ValueString = node.ValueString
	ValueInt    = node.ValueInt
	ValueBool   = node.ValueBool
	ValueFloat  = node.ValueFloat
)

type (
	NodeKind  = node.NodeKind
	ValueType = node.ValueType
)

var (
	Extract        = node.Extract
	FromTags       = node.FromTags
	ParseValueType = node.ParseValueType
)

type Path = path.Path

var (
	NewPath   = path.NewPath
	ParsePath = path.ParsePath
)
type (
	Visitor      = walk.Visitor
	VisitContext = walk.VisitContext
	Walker       = walk.Walker
	WalkConfig   = walk.WalkConfig
	WalkOption   = walk.WalkOption
)

var (
	NewWalker    = walk.NewWalker
	NewContext   = walk.NewContext
	WithMaxDepth = walk.WithMaxDepth
)
