package sidebar


func getValueAtPath(values map[string]any, path string) any {
	if path == "" {
		return nil
	}
	return values[path]
}
