package render

import (
	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui/components/containers"
)

func extractSummary(node schema.Node, value any, maxFields int) containers.Summary {
	return containers.Extract(node, value, maxFields)
}

func formatSummary(s containers.Summary) string {
	return containers.Format(s)
}
