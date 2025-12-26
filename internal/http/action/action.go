package action

import (
	"net/url"
	"strconv"
	"strings"
)

type ActionType string

const (
	ActionSave    ActionType = "save"
	ActionAdd     ActionType = "add"
	ActionRemove  ActionType = "remove"
	ActionConfirm ActionType = "confirm"
)

type Action struct {
	Type  ActionType
	Field string
	Index int
}

func Parse(form url.Values) Action {
	value := form.Get("action")
	if value == "" {
		return Action{Type: ActionSave}
	}

	parts := strings.Split(value, ":")

	switch parts[0] {
	case "add":
		if len(parts) < 2 || parts[1] == "" {
			return Action{Type: ActionSave}
		}
		return Action{
			Type:  ActionAdd,
			Field: parts[1],
		}

	case "remove":
		if len(parts) < 3 || parts[1] == "" || parts[2] == "" {
			return Action{Type: ActionSave}
		}
		index, err := strconv.Atoi(parts[2])
		if err != nil {
			return Action{Type: ActionSave}
		}
		return Action{
			Type:  ActionRemove,
			Field: parts[1],
			Index: index,
		}

	case "confirm":
		return Action{Type: ActionConfirm}

	default:
		return Action{Type: ActionSave}
	}
}
