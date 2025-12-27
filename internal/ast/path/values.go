package path

// ValuesByPath maps dotted field paths to their current values.
// Keys are path strings like "Database.Host" or "Services.0.Port".
// Values are the actual field values (string, int, bool, etc.).
type ValuesByPath = map[string]any
