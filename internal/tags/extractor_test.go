package tags

import "testing"

func TestExtract_SingleField(t *testing.T) {
	type Config struct {
		Host string `circuit:"text,help:Server hostname"`
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	f := fields[0]
	if f.Name != "Host" {
		t.Errorf("expected name Host, got %s", f.Name)
	}
	if f.Type != "string" {
		t.Errorf("expected type string, got %s", f.Type)
	}
	if f.InputType != "text" {
		t.Errorf("expected input type text, got %s", f.InputType)
	}
	if f.Help != "Server hostname" {
		t.Errorf("expected help 'Server hostname', got %s", f.Help)
	}
	if f.Required {
		t.Error("expected required to be false")
	}
}

func TestExtract_MultipleFields(t *testing.T) {
	type Config struct {
		Host string `circuit:"text,help:Server hostname,required"`
		Port int    `circuit:"number,help:Server port"`
		TLS  bool   `circuit:"checkbox,help:Enable TLS"`
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	host := fields[0]
	if host.Name != "Host" {
		t.Errorf("expected name Host, got %s", host.Name)
	}
	if host.Type != "string" {
		t.Errorf("expected type string, got %s", host.Type)
	}
	if host.InputType != "text" {
		t.Errorf("expected input type text, got %s", host.InputType)
	}
	if !host.Required {
		t.Error("expected required to be true")
	}

	port := fields[1]
	if port.Name != "Port" {
		t.Errorf("expected name Port, got %s", port.Name)
	}
	if port.Type != "int" {
		t.Errorf("expected type int, got %s", port.Type)
	}
	if port.InputType != "number" {
		t.Errorf("expected input type number, got %s", port.InputType)
	}

	tls := fields[2]
	if tls.Name != "TLS" {
		t.Errorf("expected name TLS, got %s", tls.Name)
	}
	if tls.Type != "bool" {
		t.Errorf("expected type bool, got %s", tls.Type)
	}
	if tls.InputType != "checkbox" {
		t.Errorf("expected input type checkbox, got %s", tls.InputType)
	}
}

func TestExtract_IgnoredTag(t *testing.T) {
	type Config struct {
		Host string `circuit:"-"`
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(fields))
	}
}

func TestExtract_ZeroConfig(t *testing.T) {
	type Config struct {
		Host string
		Port int
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	if fields[0].Name != "Host" || fields[0].InputType != "text" {
		t.Errorf("expected Host/text, got %s/%s", fields[0].Name, fields[0].InputType)
	}
	if fields[1].Name != "Port" || fields[1].InputType != "number" {
		t.Errorf("expected Port/number, got %s/%s", fields[1].Name, fields[1].InputType)
	}
}

func TestExtract_NonPointer(t *testing.T) {
	type Config struct {
		Host string `circuit:"text"`
	}

	cfg := Config{}
	_, err := Extract(cfg)
	if err == nil {
		t.Fatal("expected error when passing non-pointer")
	}
}

func TestExtract_NonStruct(t *testing.T) {
	s := "test"
	_, err := Extract(&s)
	if err == nil {
		t.Fatal("expected error when passing non-struct")
	}
}

func TestExtract_AdvancedTags(t *testing.T) {
	type Config struct {
		Level    int    `circuit:"type:range,min:0,max:100,step:5"`
		Category string `circuit:"type:select,options:A=Option A;B=Option B"`
		Gender   string `circuit:"type:radio,options:M=Male;F=Female"`
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	level := fields[0]
	if level.InputType != "range" {
		t.Errorf("expected input type range, got %s", level.InputType)
	}
	if level.Min != "0" {
		t.Errorf("expected min 0, got %s", level.Min)
	}
	if level.Max != "100" {
		t.Errorf("expected max 100, got %s", level.Max)
	}
	if level.Step != "5" {
		t.Errorf("expected step 5, got %s", level.Step)
	}

	category := fields[1]
	if category.InputType != "select" {
		t.Errorf("expected input type select, got %s", category.InputType)
	}
	if len(category.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(category.Options))
	}
	if category.Options[0].Label != "A" || category.Options[0].Value != "Option A" {
		t.Errorf("expected option A=Option A, got %v", category.Options[0])
	}

	gender := fields[2]
	if gender.InputType != "radio" {
		t.Errorf("expected input type radio, got %s", gender.InputType)
	}
	if len(gender.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(gender.Options))
	}
	if gender.Options[1].Label != "F" || gender.Options[1].Value != "Female" {
		t.Errorf("expected option F=Female, got %v", gender.Options[1])
	}
}

func TestExtract_ComplexTags(t *testing.T) {
	type Config struct {
		Password string `circuit:"type:password,required,help:Enter your password"`
		Count    int    `circuit:"type:number,help:Item count"`
		Active   bool   `circuit:"checkbox"`
		Email    string `circuit:"email,required"`
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}

	pass := fields[0]
	if pass.Name != "Password" {
		t.Errorf("expected name Password, got %s", pass.Name)
	}
	if pass.InputType != "password" {
		t.Errorf("expected input type password, got %s", pass.InputType)
	}
	if !pass.Required {
		t.Error("expected required to be true")
	}
	if pass.Help != "Enter your password" {
		t.Errorf("expected help 'Enter your password', got %s", pass.Help)
	}

	count := fields[1]
	if count.Name != "Count" {
		t.Errorf("expected name Count, got %s", count.Name)
	}
	if count.InputType != "number" {
		t.Errorf("expected input type number, got %s", count.InputType)
	}
	if count.Help != "Item count" {
		t.Errorf("expected help 'Item count', got %s", count.Help)
	}

	active := fields[2]
	if active.Name != "Active" {
		t.Errorf("expected name Active, got %s", active.Name)
	}
	if active.InputType != "checkbox" {
		t.Errorf("expected input type checkbox, got %s", active.InputType)
	}

	email := fields[3]
	if email.Name != "Email" {
		t.Errorf("expected name Email, got %s", email.Name)
	}
	if email.InputType != "email" {
		t.Errorf("expected input type email, got %s", email.InputType)
	}
	if !email.Required {
		t.Error("expected required to be true")
	}
}

func TestExtract_PointerFields(t *testing.T) {
	type Config struct {
		Host *string
		Port *int
		Flag *bool
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	if fields[0].Type != "string" {
		t.Errorf("expected type string for *string, got %s", fields[0].Type)
	}
	if fields[0].InputType != "text" {
		t.Errorf("expected input type text for *string, got %s", fields[0].InputType)
	}

	if fields[1].Type != "int" {
		t.Errorf("expected type int for *int, got %s", fields[1].Type)
	}
	if fields[1].InputType != "number" {
		t.Errorf("expected input type number for *int, got %s", fields[1].InputType)
	}

	if fields[2].Type != "bool" {
		t.Errorf("expected type bool for *bool, got %s", fields[2].Type)
	}
	if fields[2].InputType != "checkbox" {
		t.Errorf("expected input type checkbox for *bool, got %s", fields[2].InputType)
	}
}
