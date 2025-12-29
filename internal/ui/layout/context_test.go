package layout

import (
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ui/render"
)

func TestNewPageContext(t *testing.T) {
	schema := &ast.Schema{Name: "AppConfig"}
	values := ast.ValuesByPath{"host": "localhost"}

	rc := render.NewRenderContext(schema, values)
	pc := NewPageContext(rc)

	if pc.RenderContext != rc {
		t.Error("RenderContext not embedded correctly")
	}

	expectedTitle := "AppConfig Configuration"
	if pc.Title != expectedTitle {
		t.Errorf("Expected title %q, got %q", expectedTitle, pc.Title)
	}

	if pc.Brand {
		t.Error("Brand should default to false")
	}

	if pc.TopContent != nil {
		t.Error("TopContent should default to nil")
	}
}

func TestPageContextEmbedding(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := render.NewRenderContext(schema, values)
	rc.ReadOnly = true
	rc.MaxDepth = 3

	pc := NewPageContext(rc)

	// Verify embedded fields are accessible
	if !pc.ReadOnly {
		t.Error("ReadOnly should be accessible via embedding")
	}

	if pc.MaxDepth != 3 {
		t.Errorf("Expected MaxDepth=3, got %d", pc.MaxDepth)
	}

	if pc.Schema != schema {
		t.Error("Schema should be accessible via embedding")
	}
}

func TestPageContextWithTopContent(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := render.NewRenderContext(schema, values)
	pc := NewPageContext(rc)

	topBanner := h.Div(h.Class("banner"), g.Text("Warning"))
	pc.TopContent = []g.Node{topBanner}

	if len(pc.TopContent) != 1 {
		t.Errorf("Expected 1 TopContent node, got %d", len(pc.TopContent))
	}
}

func TestPageContextWithCustomTitle(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := render.NewRenderContext(schema, values)
	pc := NewPageContext(rc)

	customTitle := "My Custom Dashboard"
	pc.Title = customTitle

	if pc.Title != customTitle {
		t.Errorf("Expected title %q, got %q", customTitle, pc.Title)
	}
}

func TestPageContextWithBrand(t *testing.T) {
	schema := &ast.Schema{Name: "Test"}
	values := ast.ValuesByPath{}

	rc := render.NewRenderContext(schema, values)
	pc := NewPageContext(rc)

	pc.Brand = true

	if !pc.Brand {
		t.Error("Brand should be true")
	}
}
