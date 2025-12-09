package form

import (
	"net/url"
	"testing"

	"github.com/moq77111113/circuit/internal/schema"
)

type ConvertConfig struct {
	Name  string `yaml:"name" circuit:"type:text"`
	Port  int    `yaml:"port" circuit:"type:number"`
	Debug bool   `yaml:"debug" circuit:"type:checkbox"`
}

func TestExtractValues(t *testing.T) {
	cfg := ConvertConfig{
		Name:  "test",
		Port:  8080,
		Debug: true,
	}

	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	values := ExtractValues(&cfg, s)

	if values["Name"] != "test" {
		t.Errorf("expected Name=test, got %v", values["Name"])
	}
	if values["Port"] != 8080 {
		t.Errorf("expected Port=8080, got %v", values["Port"])
	}
	if values["Debug"] != true {
		t.Errorf("expected Debug=true, got %v", values["Debug"])
	}
}

func TestApplyForm_String(t *testing.T) {
	cfg := ConvertConfig{}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Name", "newname")

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "newname" {
		t.Errorf("expected Name=newname, got %s", cfg.Name)
	}
}

func TestApplyForm_Int(t *testing.T) {
	cfg := ConvertConfig{}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Port", "9000")

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Port != 9000 {
		t.Errorf("expected Port=9000, got %d", cfg.Port)
	}
}

func TestApplyForm_IntEmpty(t *testing.T) {
	cfg := ConvertConfig{Port: 8080}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Port", "")

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Port != 0 {
		t.Errorf("expected Port=0 for empty value, got %d", cfg.Port)
	}
}

func TestApplyForm_IntInvalid(t *testing.T) {
	cfg := ConvertConfig{}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Port", "invalid")

	err = Apply(&cfg, s, form)
	if err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestApplyForm_Bool(t *testing.T) {
	cfg := ConvertConfig{}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Debug", "on")

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Debug {
		t.Error("expected Debug=true when value is 'on'")
	}
}

func TestApplyForm_BoolOff(t *testing.T) {
	cfg := ConvertConfig{Debug: true}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	// Debug not set = checkbox unchecked

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Debug {
		t.Error("expected Debug=false when not in form")
	}
}

func TestApplyForm_Multiple(t *testing.T) {
	cfg := ConvertConfig{}
	s, err := schema.Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Set("Name", "app")
	form.Set("Port", "3000")
	form.Set("Debug", "on")

	err = Apply(&cfg, s, form)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "app" {
		t.Errorf("expected Name=app, got %s", cfg.Name)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected Port=3000, got %d", cfg.Port)
	}
	if !cfg.Debug {
		t.Errorf("expected Debug=true, got %v", cfg.Debug)
	}
}
