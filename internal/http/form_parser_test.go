package http

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseIndexedField_StringSlice(t *testing.T) {
	form := url.Values{
		"tags.0": {"golang"},
		"tags.1": {"rust"},
		"tags.2": {"python"},
	}

	result := parseIndexedField(form, "tags")

	expected := []string{"golang", "rust", "python"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseIndexedField_WithGaps(t *testing.T) {
	form := url.Values{
		"tags.0": {"golang"},
		"tags.2": {"python"},
		"tags.5": {"java"},
	}

	result := parseIndexedField(form, "tags")

	expected := []string{"golang", "python", "java"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseIndexedField_OutOfOrder(t *testing.T) {
	form := url.Values{
		"tags.2": {"python"},
		"tags.0": {"golang"},
		"tags.1": {"rust"},
	}

	result := parseIndexedField(form, "tags")

	expected := []string{"golang", "rust", "python"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseIndexedField_Empty(t *testing.T) {
	form := url.Values{}

	result := parseIndexedField(form, "tags")

	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestParseIndexedField_NonIndexedField(t *testing.T) {
	form := url.Values{
		"host": {"localhost"},
		"port": {"8080"},
	}

	result := parseIndexedField(form, "tags")

	if len(result) != 0 {
		t.Errorf("expected empty slice for non-existent field, got %v", result)
	}
}

func TestParseIndexedField_MixedFields(t *testing.T) {
	form := url.Values{
		"tags.0":  {"go"},
		"tags.1":  {"rust"},
		"host":    {"localhost"},
		"ports.0": {"8080"},
		"ports.1": {"9090"},
	}

	tagsResult := parseIndexedField(form, "tags")
	portsResult := parseIndexedField(form, "ports")

	expectedTags := []string{"go", "rust"}
	expectedPorts := []string{"8080", "9090"}

	if !reflect.DeepEqual(tagsResult, expectedTags) {
		t.Errorf("expected tags %v, got %v", expectedTags, tagsResult)
	}

	if !reflect.DeepEqual(portsResult, expectedPorts) {
		t.Errorf("expected ports %v, got %v", expectedPorts, portsResult)
	}
}

func TestParseIndexedField_EmptyValues(t *testing.T) {
	form := url.Values{
		"tags.0": {""},
		"tags.1": {"rust"},
		"tags.2": {""},
	}

	result := parseIndexedField(form, "tags")

	expected := []string{"", "rust", ""}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
