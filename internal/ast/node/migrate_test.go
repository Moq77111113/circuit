package node

import (
	"testing"

	"github.com/moq77111113/circuit/internal/tags"
)

func TestFromTags_Primitives(t *testing.T) {
	fields := []tags.Field{
		{Name: "Host", Type: "string", InputType: tags.TypeText},
		{Name: "Port", Type: "int", InputType: tags.TypeNumber},
		{Name: "Enabled", Type: "bool", InputType: tags.TypeCheckbox},
	}

	nodes := FromTags(fields)

	if len(nodes) != 3 {
		t.Fatalf("len(nodes) = %d, want 3", len(nodes))
	}

	if nodes[0].Kind != KindPrimitive {
		t.Errorf("nodes[0].Kind = %v, want KindPrimitive", nodes[0].Kind)
	}
	if nodes[0].ValueType != ValueString {
		t.Errorf("nodes[0].ValueType = %v, want ValueString", nodes[0].ValueType)
	}

	if nodes[1].Kind != KindPrimitive {
		t.Errorf("nodes[1].Kind = %v, want KindPrimitive", nodes[1].Kind)
	}
	if nodes[1].ValueType != ValueInt {
		t.Errorf("nodes[1].ValueType = %v, want ValueInt", nodes[1].ValueType)
	}

	if nodes[2].Kind != KindPrimitive {
		t.Errorf("nodes[2].Kind = %v, want KindPrimitive", nodes[2].Kind)
	}
	if nodes[2].ValueType != ValueBool {
		t.Errorf("nodes[2].ValueType = %v, want ValueBool", nodes[2].ValueType)
	}
}

func TestFromTags_NestedStruct(t *testing.T) {
	fields := []tags.Field{
		{
			Name:      "Database",
			InputType: tags.TypeSection,
			Fields: []tags.Field{
				{Name: "Host", Type: "string", InputType: tags.TypeText},
				{Name: "Port", Type: "int", InputType: tags.TypeNumber},
			},
		},
	}

	nodes := FromTags(fields)

	if len(nodes) != 1 {
		t.Fatalf("len(nodes) = %d, want 1", len(nodes))
	}

	db := nodes[0]
	if db.Kind != KindStruct {
		t.Errorf("db.Kind = %v, want KindStruct", db.Kind)
	}
	if db.Name != "Database" {
		t.Errorf("db.Name = %s, want Database", db.Name)
	}
	if len(db.Children) != 2 {
		t.Fatalf("len(db.Children) = %d, want 2", len(db.Children))
	}

	if db.Children[0].Kind != KindPrimitive {
		t.Errorf("child[0].Kind = %v, want KindPrimitive", db.Children[0].Kind)
	}
	if db.Children[0].Name != "Host" {
		t.Errorf("child[0].Name = %s, want Host", db.Children[0].Name)
	}
}

func TestFromTags_PrimitiveSlice(t *testing.T) {
	fields := []tags.Field{
		{
			Name:        "Tags",
			IsSlice:     true,
			Type:        "slice",
			ElementType: "string",
			InputType:   tags.TypeText,
		},
		{
			Name:        "Ports",
			IsSlice:     true,
			Type:        "slice",
			ElementType: "int",
			InputType:   tags.TypeNumber,
		},
	}

	nodes := FromTags(fields)

	if len(nodes) != 2 {
		t.Fatalf("len(nodes) = %d, want 2", len(nodes))
	}

	tags := nodes[0]
	if tags.Kind != KindSlice {
		t.Errorf("tags.Kind = %v, want KindSlice", tags.Kind)
	}
	if tags.ElementKind != KindPrimitive {
		t.Errorf("tags.ElementKind = %v, want KindPrimitive", tags.ElementKind)
	}
	if tags.ValueType != ValueString {
		t.Errorf("tags.ValueType = %v, want ValueString", tags.ValueType)
	}

	ports := nodes[1]
	if ports.ElementKind != KindPrimitive {
		t.Errorf("ports.ElementKind = %v, want KindPrimitive", ports.ElementKind)
	}
	if ports.ValueType != ValueInt {
		t.Errorf("ports.ValueType = %v, want ValueInt", ports.ValueType)
	}
}

func TestFromTags_StructSlice(t *testing.T) {
	fields := []tags.Field{
		{
			Name:        "Services",
			IsSlice:     true,
			Type:        "slice",
			ElementType: "struct",
			Fields: []tags.Field{
				{Name: "Name", Type: "string", InputType: tags.TypeText},
				{Name: "Port", Type: "int", InputType: tags.TypeNumber},
			},
		},
	}

	nodes := FromTags(fields)

	if len(nodes) != 1 {
		t.Fatalf("len(nodes) = %d, want 1", len(nodes))
	}

	services := nodes[0]
	if services.Kind != KindSlice {
		t.Errorf("services.Kind = %v, want KindSlice", services.Kind)
	}
	if services.ElementKind != KindStruct {
		t.Errorf("services.ElementKind = %v, want KindStruct", services.ElementKind)
	}
	if len(services.Children) != 2 {
		t.Fatalf("len(services.Children) = %d, want 2", len(services.Children))
	}

	if services.Children[0].Name != "Name" {
		t.Errorf("child[0].Name = %s, want Name", services.Children[0].Name)
	}
	if services.Children[0].Kind != KindPrimitive {
		t.Errorf("child[0].Kind = %v, want KindPrimitive", services.Children[0].Kind)
	}
}

func TestFromTags_DeepNested(t *testing.T) {
	fields := []tags.Field{
		{
			Name:        "Services",
			IsSlice:     true,
			Type:        "slice",
			ElementType: "struct",
			Fields: []tags.Field{
				{Name: "Name", Type: "string", InputType: tags.TypeText},
				{
					Name:        "Endpoints",
					IsSlice:     true,
					Type:        "slice",
					ElementType: "struct",
					Fields: []tags.Field{
						{Name: "Path", Type: "string", InputType: tags.TypeText},
						{
							Name:        "AllowedRoles",
							IsSlice:     true,
							Type:        "slice",
							ElementType: "string",
							InputType:   tags.TypeText,
						},
					},
				},
			},
		},
	}

	nodes := FromTags(fields)

	if len(nodes) != 1 {
		t.Fatalf("len(nodes) = %d, want 1", len(nodes))
	}

	services := nodes[0]
	if services.Kind != KindSlice || services.ElementKind != KindStruct {
		t.Fatal("services should be KindSlice of KindStruct")
	}

	endpoints := services.Children[1]
	if endpoints.Name != "Endpoints" {
		t.Errorf("endpoints.Name = %s, want Endpoints", endpoints.Name)
	}
	if endpoints.Kind != KindSlice || endpoints.ElementKind != KindStruct {
		t.Fatal("endpoints should be KindSlice of KindStruct")
	}

	roles := endpoints.Children[1]
	if roles.Name != "AllowedRoles" {
		t.Errorf("roles.Name = %s, want AllowedRoles", roles.Name)
	}
	if roles.Kind != KindSlice || roles.ElementKind != KindPrimitive {
		t.Fatal("roles should be KindSlice of KindPrimitive")
	}
	if roles.ValueType != ValueString {
		t.Errorf("roles.ValueType = %v, want ValueString", roles.ValueType)
	}
}

func TestFromTags_PreservesMetadata(t *testing.T) {
	fields := []tags.Field{
		{
			Name:      "Port",
			Type:      "int",
			InputType: tags.TypeNumber,
			Help:      "Server port",
			Required:  true,
			Min:       "1",
			Max:       "65535",
			Step:      "1",
		},
	}

	nodes := FromTags(fields)

	if len(nodes) != 1 {
		t.Fatalf("len(nodes) = %d, want 1", len(nodes))
	}

	node := nodes[0]
	if node.UI.Help != "Server port" {
		t.Errorf("node.UI.Help = %s, want 'Server port'", node.UI.Help)
	}
	if !node.UI.Required {
		t.Error("node.UI.Required should be true")
	}
	if node.UI.Min != "1" {
		t.Errorf("node.UI.Min = %s, want '1'", node.UI.Min)
	}
	if node.UI.Max != "65535" {
		t.Errorf("node.UI.Max = %s, want '65535'", node.UI.Max)
	}
	if node.UI.Step != "1" {
		t.Errorf("node.UI.Step = %s, want '1'", node.UI.Step)
	}
}
