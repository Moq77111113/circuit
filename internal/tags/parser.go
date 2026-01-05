package tags

import (
	"strconv"
	"strings"
)

var tagHandlers = map[string]func(*Field, string){
	"type":    func(f *Field, v string) { f.InputType = InputType(v) },
	"help":    func(f *Field, v string) { f.Help = v },
	"min":     func(f *Field, v string) { f.Min = v },
	"max":     func(f *Field, v string) { f.Max = v },
	"step":    func(f *Field, v string) { f.Step = v },
	"pattern": func(f *Field, v string) { f.Pattern = v },
	"minlen": func(f *Field, v string) {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 0 {
			f.MinLen = n
		}
	},
	"maxlen": func(f *Field, v string) {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 0 {
			f.MaxLen = n
		}
	},
	"options": func(f *Field, v string) {
		opts := strings.SplitSeq(v, ";")
		for opt := range opts {
			kv := strings.SplitN(opt, "=", 2)
			if len(kv) == 2 {
				f.Options = append(f.Options, Option{
					Value: strings.TrimSpace(kv[0]),
					Label: strings.TrimSpace(kv[1]),
				})
			} else {
				val := strings.TrimSpace(kv[0])
				f.Options = append(f.Options, Option{
					Label: val,
					Value: val,
				})
			}
		}
	},
}

var flagActions = map[string]func(*Field){
	"required": func(f *Field) { f.Required = true },
	"readonly": func(f *Field) { f.ReadOnly = true },
}

func parseTag(tag string, f *Field) {
	parts := strings.Split(tag, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			if handler, ok := tagHandlers[key]; ok {
				handler(f, val)
			}
			continue
		}

		if handler, ok := flagActions[part]; ok {
			handler(f)
		} else if i == 0 {
			f.InputType = InputType(part)
		}
	}
}
