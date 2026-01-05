package validation

// presetPatterns maps preset names to regex patterns.
var presetPatterns = map[string]string{
	"email": `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	"url":   `^https?://[^\s/$.?#].[^\s]*$`,
	"phone": `^\+?[0-9\s\-\(\)]{7,20}$`,
}

// IsPreset checks if a pattern name is a known preset.
func IsPreset(name string) bool {
	_, ok := presetPatterns[name]
	return ok
}

// GetPreset retrieves the regex for a preset name.
func GetPreset(name string) (string, bool) {
	pattern, ok := presetPatterns[name]
	return pattern, ok
}
