package schema

import (
	"testing"

	"github.com/moq77111113/circuit/internal/tags"
)

func TestNode_PrimitiveConstruction(t *testing.T) {
	tests := []struct {
		name      string
		node      Node
		wantKind  NodeKind
		wantValue ValueType
	}{
		{
			name: "string primitive",
			node: Node{
				Name:      "Host",
				Kind:      KindPrimitive,
				ValueType: ValueString,
				InputType: tags.TypeText,
			},
			wantKind:  KindPrimitive,
			wantValue: ValueString,
		},
		{
			name: "int primitive",
			node: Node{
				Name:      "Port",
				Kind:      KindPrimitive,
				ValueType: ValueInt,
				InputType: tags.TypeNumber,
			},
			wantKind:  KindPrimitive,
			wantValue: ValueInt,
		},
		{
			name: "bool primitive",
			node: Node{
				Name:      "Enabled",
				Kind:      KindPrimitive,
				ValueType: ValueBool,
				InputType: tags.TypeCheckbox,
			},
			wantKind:  KindPrimitive,
			wantValue: ValueBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.node.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", tt.node.Kind, tt.wantKind)
			}
			if tt.node.ValueType != tt.wantValue {
				t.Errorf("ValueType = %v, want %v", tt.node.ValueType, tt.wantValue)
			}
		})
	}
}

func TestNode_StructConstruction(t *testing.T) {
	node := Node{
		Name: "Database",
		Kind: KindStruct,
		Children: []Node{
			{
				Name:      "Host",
				Kind:      KindPrimitive,
				ValueType: ValueString,
				InputType: tags.TypeText,
			},
			{
				Name:      "Port",
				Kind:      KindPrimitive,
				ValueType: ValueInt,
				InputType: tags.TypeNumber,
			},
		},
	}

	if node.Kind != KindStruct {
		t.Errorf("Kind = %v, want %v", node.Kind, KindStruct)
	}
	if len(node.Children) != 2 {
		t.Fatalf("len(Children) = %d, want 2", len(node.Children))
	}
	if node.Children[0].Name != "Host" {
		t.Errorf("Children[0].Name = %s, want Host", node.Children[0].Name)
	}
}

func TestNode_SliceConstruction(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		wantKind    NodeKind
		wantElement NodeKind
	}{
		{
			name: "primitive slice ([]string)",
			node: Node{
				Name:        "Tags",
				Kind:        KindSlice,
				ElementKind: KindPrimitive,
				ValueType:   ValueString,
				InputType:   tags.TypeText,
			},
			wantKind:    KindSlice,
			wantElement: KindPrimitive,
		},
		{
			name: "struct slice ([]Service)",
			node: Node{
				Name:        "Services",
				Kind:        KindSlice,
				ElementKind: KindStruct,
				Children: []Node{
					{
						Name:      "Name",
						Kind:      KindPrimitive,
						ValueType: ValueString,
						InputType: tags.TypeText,
					},
				},
			},
			wantKind:    KindSlice,
			wantElement: KindStruct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.node.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", tt.node.Kind, tt.wantKind)
			}
			if tt.node.ElementKind != tt.wantElement {
				t.Errorf("ElementKind = %v, want %v", tt.node.ElementKind, tt.wantElement)
			}
		})
	}
}

func TestNode_DeepNested(t *testing.T) {
	// Services[].Endpoints[].AllowedRoles[]
	node := Node{
		Name:        "Services",
		Kind:        KindSlice,
		ElementKind: KindStruct,
		Children: []Node{
			{
				Name:      "Name",
				Kind:      KindPrimitive,
				ValueType: ValueString,
				InputType: tags.TypeText,
			},
			{
				Name:        "Endpoints",
				Kind:        KindSlice,
				ElementKind: KindStruct,
				Children: []Node{
					{
						Name:      "Path",
						Kind:      KindPrimitive,
						ValueType: ValueString,
						InputType: tags.TypeText,
					},
					{
						Name:        "AllowedRoles",
						Kind:        KindSlice,
						ElementKind: KindPrimitive,
						ValueType:   ValueString,
						InputType:   tags.TypeText,
					},
				},
			},
		},
	}

	if node.Kind != KindSlice {
		t.Fatalf("root Kind = %v, want KindSlice", node.Kind)
	}
	if node.ElementKind != KindStruct {
		t.Fatalf("root ElementKind = %v, want KindStruct", node.ElementKind)
	}

	endpoints := node.Children[1]
	if endpoints.Name != "Endpoints" {
		t.Errorf("endpoints.Name = %s, want Endpoints", endpoints.Name)
	}
	if endpoints.Kind != KindSlice {
		t.Errorf("endpoints.Kind = %v, want KindSlice", endpoints.Kind)
	}
	if endpoints.ElementKind != KindStruct {
		t.Errorf("endpoints.ElementKind = %v, want KindStruct", endpoints.ElementKind)
	}

	roles := endpoints.Children[1]
	if roles.Name != "AllowedRoles" {
		t.Errorf("roles.Name = %s, want AllowedRoles", roles.Name)
	}
	if roles.Kind != KindSlice {
		t.Errorf("roles.Kind = %v, want KindSlice", roles.Kind)
	}
	if roles.ElementKind != KindPrimitive {
		t.Errorf("roles.ElementKind = %v, want KindPrimitive", roles.ElementKind)
	}
}
