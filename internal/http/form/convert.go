package form

import (
	"net/url"

	"github.com/moq77111113/circuit/internal/ast"
)

// Apply updates a config struct from form data using the Visitor pattern.
func Apply(cfg any, s ast.Schema, form url.Values) error {
	return ApplyNodes(cfg, s.Nodes, form)
}
