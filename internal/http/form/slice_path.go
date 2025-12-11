package form

// pathSegment represents a parsed path segment with optional index.
type pathSegment struct {
	name     string
	hasIndex bool
	index    int
}

// extractPathSegments parses a path into segments with indices.
// "Services.0.Endpoints" -> [{Services, true, 0}, {Endpoints, false, 0}]
// "Database.Host" -> [{Database, false, 0}, {Host, false, 0}]
func extractPathSegments(path string) []pathSegment {
	if path == "" {
		return nil
	}

	var segments []pathSegment
	parts := parseFieldNames(path)

	pathParts := []string{}
	current := ""
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			if current != "" {
				pathParts = append(pathParts, current)
				current = ""
			}
		} else {
			current += string(path[i])
		}
	}
	if current != "" {
		pathParts = append(pathParts, current)
	}

	fieldIdx := 0
	for i := 0; i < len(pathParts); i++ {
		part := pathParts[i]

		isNumber := true
		for _, c := range part {
			if c < '0' || c > '9' {
				isNumber = false
				break
			}
		}

		if isNumber {
			if len(segments) > 0 {
				idx := 0
				for _, c := range part {
					idx = idx*10 + int(c-'0')
				}
				segments[len(segments)-1].hasIndex = true
				segments[len(segments)-1].index = idx
			}
		} else {
			if fieldIdx < len(parts) && parts[fieldIdx] == part {
				segments = append(segments, pathSegment{name: part, hasIndex: false, index: 0})
				fieldIdx++
			}
		}
	}

	return segments
}

// parseFieldNames extracts field names from a path, ignoring indices.
// "Services.0.Endpoints" -> ["Services", "Endpoints"]
// "Database.Host" -> ["Database", "Host"]
func parseFieldNames(path string) []string {
	if path == "" {
		return nil
	}

	var names []string
	current := ""

	for i := 0; i < len(path); i++ {
		ch := path[i]
		if ch == '.' {
			if current != "" {
				isNumber := true
				for _, c := range current {
					if c < '0' || c > '9' {
						isNumber = false
						break
					}
				}
				if !isNumber {
					names = append(names, current)
				}
				current = ""
			}
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		isNumber := true
		for _, c := range current {
			if c < '0' || c > '9' {
				isNumber = false
				break
			}
		}
		if !isNumber {
			names = append(names, current)
		}
	}

	return names
}
