package inputs

import (
	"fmt"
	"strconv"

	"github.com/moq77111113/circuit/internal/tags"
	"github.com/moq77111113/circuit/internal/ui/styles"
	"github.com/moq77111113/circuit/internal/validation"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Text(field tags.Field, value any) g.Node {
	return simpleInput(field, value, "text")
}

func Date(field tags.Field, value any) g.Node {
	return simpleInput(field, value, "date")
}

func Time(field tags.Field, value any) g.Node {
	return simpleInput(field, value, "time")
}

func Password(field tags.Field, value any) g.Node {
	return simpleInput(field, value, "password")
}

func simpleInput(field tags.Field, value any, inputType string) g.Node {
	attrs := BaseAttrs(field)
	if value != nil {
		attrs = append(attrs, h.Value(fmt.Sprintf("%v", value)))
	}

	// HTML5 validation attributes
	if field.Pattern != "" {
		pattern := field.Pattern
		if preset, ok := validation.GetPreset(field.Pattern); ok {
			pattern = preset
		}
		attrs = append(attrs, g.Attr("pattern", pattern))
	}
	if field.MinLen > 0 {
		attrs = append(attrs, g.Attr("minlength", strconv.Itoa(field.MinLen)))
	}
	if field.MaxLen > 0 {
		attrs = append(attrs, g.Attr("maxlength", strconv.Itoa(field.MaxLen)))
	}

	return h.Input(append(attrs, h.Type(inputType))...)
}

func Number(field tags.Field, value any) g.Node {
	attrs := BaseAttrs(field)
	if value != nil {
		attrs = append(attrs, h.Value(fmt.Sprintf("%v", value)))
	}
	if field.Min != "" {
		attrs = append(attrs, h.Min(field.Min))
	}
	if field.Max != "" {
		attrs = append(attrs, h.Max(field.Max))
	}
	if field.Step != "" {
		attrs = append(attrs, h.Step(field.Step))
	}
	return h.Input(append(attrs, h.Type("number"))...)
}

func Range(field tags.Field, value any) g.Node {
	attrs := BaseAttrs(field)
	if value != nil {
		attrs = append(attrs, h.Value(fmt.Sprintf("%v", value)))
	}
	if field.Min != "" {
		attrs = append(attrs, h.Min(field.Min))
	}
	if field.Max != "" {
		attrs = append(attrs, h.Max(field.Max))
	}
	if field.Step != "" {
		attrs = append(attrs, h.Step(field.Step))
	}

	attrs = append(attrs, g.Attr("oninput", "this.nextElementSibling.nextElementSibling.value = this.value"))

	return h.Div(
		h.Class(styles.RangeWrapper),
		h.Span(h.Class(styles.RangeMin), g.Text(field.Min)),
		h.Input(append(attrs, h.Type("range"))...),
		h.Span(h.Class(styles.RangeMax), g.Text(field.Max)),
		g.El("output", h.Class(styles.RangeValue), g.Text(fmt.Sprintf("%v", value))),
	)
}
