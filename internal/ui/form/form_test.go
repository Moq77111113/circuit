package form

import (
	"strings"
	"testing"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

func TestForm_TextInput(t *testing.T) {
	s := schema.Schema{
		Name: "Config",
		Fields: []tags.Field{
			{
				Name:      "Host",
				Type:      "string",
				InputType: "text",
				Help:      "Server hostname",
				Required:  false,
			},
		},
	}

	html := renderToString(Form(s, nil))

	if !strings.Contains(html, `type="text"`) {
		t.Error("expected text input type")
	}
	if !strings.Contains(html, `name="Host"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "Server hostname") {
		t.Error("expected help text")
	}
}

func TestForm_NumberInput(t *testing.T) {
	s := schema.Schema{
		Name: "Config",
		Fields: []tags.Field{
			{
				Name:      "Port",
				Type:      "int",
				InputType: "number",
				Help:      "Server port",
				Required:  true,
			},
		},
	}

	html := renderToString(Form(s, nil))

	if !strings.Contains(html, `type="number"`) {
		t.Error("expected number input type")
	}
	if !strings.Contains(html, `name="Port"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "Server port") {
		t.Error("expected help text")
	}
	if !strings.Contains(html, "required") {
		t.Error("expected required attribute")
	}
}

func TestForm_Checkbox(t *testing.T) {
	s := schema.Schema{
		Name: "Config",
		Fields: []tags.Field{
			{
				Name:      "TLS",
				Type:      "bool",
				InputType: "checkbox",
				Help:      "Enable TLS",
				Required:  false,
			},
		},
	}

	html := renderToString(Form(s, nil))

	if !strings.Contains(html, `type="checkbox"`) {
		t.Error("expected checkbox input type")
	}
	if !strings.Contains(html, `name="TLS"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "Enable TLS") {
		t.Error("expected help text")
	}
}

func TestForm_MultipleFields(t *testing.T) {
	s := schema.Schema{
		Name: "Config",
		Fields: []tags.Field{
			{Name: "Host", Type: "string", InputType: "text", Help: "Hostname"},
			{Name: "Port", Type: "int", InputType: "number", Help: "Port"},
			{Name: "TLS", Type: "bool", InputType: "checkbox", Help: "TLS"},
		},
	}

	html := renderToString(Form(s, nil))

	if !strings.Contains(html, `name="Host"`) {
		t.Error("expected Host field")
	}
	if !strings.Contains(html, `name="Port"`) {
		t.Error("expected Port field")
	}
	if !strings.Contains(html, `name="TLS"`) {
		t.Error("expected TLS field")
	}
}

func TestForm_WithValues(t *testing.T) {
	s := schema.Schema{
		Name: "Config",
		Fields: []tags.Field{
			{Name: "Host", Type: "string", InputType: "text"},
			{Name: "Port", Type: "int", InputType: "number"},
		},
	}

	values := map[string]any{
		"Host": "localhost",
		"Port": 8080,
	}

	html := renderToString(Form(s, values))

	if !strings.Contains(html, `value="localhost"`) {
		t.Error("expected Host value")
	}
	if !strings.Contains(html, `value="8080"`) {
		t.Error("expected Port value")
	}
}

func TestForm_EmptySchema(t *testing.T) {
	s := schema.Schema{
		Name:   "Empty",
		Fields: []tags.Field{},
	}

	html := renderToString(Form(s, nil))

	if !strings.Contains(html, "<form") {
		t.Error("expected form element even for empty schema")
	}
}
