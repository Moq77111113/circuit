package validation

import (
	"net/url"
	"testing"

	"github.com/moq77111113/circuit/internal/ast/node"
	"github.com/moq77111113/circuit/internal/ast/path"
)

func TestMergeFormValues_SimpleFields(t *testing.T) {
	nodes := []node.Node{
		{
			Name:      "Host",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueString,
		},
		{
			Name:      "Port",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueInt,
		},
	}

	configValues := path.ValuesByPath{
		"Host": "prod.example.com",
		"Port": 8080,
	}

	form := url.Values{}
	form.Set("Host", "localhost")
	form.Set("Port", "3000")

	result := MergeFormValues(nodes, configValues, form)

	if result["Host"] != "localhost" {
		t.Errorf("expected Host 'localhost', got '%v'", result["Host"])
	}
	if result["Port"] != "3000" {
		t.Errorf("expected Port '3000', got '%v'", result["Port"])
	}
}

func TestMergeFormValues_PartialFormData(t *testing.T) {
	nodes := []node.Node{
		{
			Name:      "Host",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueString,
		},
		{
			Name:      "Port",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueInt,
		},
	}

	configValues := path.ValuesByPath{
		"Host": "prod.example.com",
		"Port": 8080,
	}

	form := url.Values{}
	form.Set("Host", "localhost")

	result := MergeFormValues(nodes, configValues, form)

	if result["Host"] != "localhost" {
		t.Errorf("expected Host 'localhost', got '%v'", result["Host"])
	}
	if result["Port"] != 8080 {
		t.Errorf("expected Port 8080 (from config), got '%v'", result["Port"])
	}
}

func TestMergeFormValues_NestedStructs(t *testing.T) {
	nodes := []node.Node{
		{
			Name: "Server",
			Kind: node.KindStruct,
			Children: []node.Node{
				{
					Name:      "Host",
					Kind:      node.KindPrimitive,
					ValueType: node.ValueString,
				},
				{
					Name:      "Port",
					Kind:      node.KindPrimitive,
					ValueType: node.ValueInt,
				},
			},
		},
	}

	configValues := path.ValuesByPath{
		"Server.Host": "prod.example.com",
		"Server.Port": 8080,
	}

	form := url.Values{}
	form.Set("Server.Host", "localhost")
	form.Set("Server.Port", "3000")

	result := MergeFormValues(nodes, configValues, form)

	if result["Server.Host"] != "localhost" {
		t.Errorf("expected Server.Host 'localhost', got '%v'", result["Server.Host"])
	}
	if result["Server.Port"] != "3000" {
		t.Errorf("expected Server.Port '3000', got '%v'", result["Server.Port"])
	}
}

func TestMergeFormValues_EmptyFormValue(t *testing.T) {
	nodes := []node.Node{
		{
			Name:      "Host",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueString,
		},
	}

	configValues := path.ValuesByPath{
		"Host": "prod.example.com",
	}

	form := url.Values{}
	form.Set("Host", "")

	result := MergeFormValues(nodes, configValues, form)

	if result["Host"] != "" {
		t.Errorf("expected Host '' (empty from form), got '%v'", result["Host"])
	}
}

func TestMergeFormValues_PreservesConfigWhenNoFormData(t *testing.T) {
	nodes := []node.Node{
		{
			Name:      "Host",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueString,
		},
		{
			Name:      "Port",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueInt,
		},
	}

	configValues := path.ValuesByPath{
		"Host": "prod.example.com",
		"Port": 8080,
	}

	form := url.Values{}

	result := MergeFormValues(nodes, configValues, form)

	if result["Host"] != "prod.example.com" {
		t.Errorf("expected Host 'prod.example.com' (from config), got '%v'", result["Host"])
	}
	if result["Port"] != 8080 {
		t.Errorf("expected Port 8080 (from config), got '%v'", result["Port"])
	}
}

func TestMergeFormValues_BooleanFields(t *testing.T) {
	nodes := []node.Node{
		{
			Name:      "Enabled",
			Kind:      node.KindPrimitive,
			ValueType: node.ValueBool,
		},
	}

	configValues := path.ValuesByPath{
		"Enabled": false,
	}

	form := url.Values{}
	form.Set("Enabled", "true")

	result := MergeFormValues(nodes, configValues, form)

	if result["Enabled"] != "true" {
		t.Errorf("expected Enabled 'true', got '%v'", result["Enabled"])
	}
}
