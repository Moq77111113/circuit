package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/styles"
)

func renderActions(actions []ActionButton) g.Node {
	if len(actions) == 0 {
		return nil
	}

	buttons := make([]g.Node, len(actions))
	for i, action := range actions {
		buttons[i] = renderActionButton(action)
	}

	return h.Div(
		h.Class(styles.ActionsSection),
		h.H2(h.Class(styles.ActionsSectionTitle), g.Text("Actions")),
		h.Div(h.Class(styles.ActionsButtons), g.Group(buttons)),
	)
}

func renderActionButton(action ActionButton) g.Node {
	buttonAttrs := []g.Node{
		h.Type("submit"),
		h.Class(styles.Button + " " + styles.ButtonSecondary),
		h.Title(action.Description),
	}

	if action.RequireConfirmation {
		onClick := "return confirmAction('" + action.Label + "', '" + action.Description + "')"
		buttonAttrs = append(buttonAttrs, g.Attr("onclick", onClick))
	}

	return h.Form(
		h.Method("post"),
		h.Input(
			h.Type("hidden"),
			h.Name("action"),
			h.Value("execute:"+action.Name),
		),
		h.Button(g.Group(buttonAttrs), g.Text(action.Label)),
	)
}

func renderErrorBanner(msg string) g.Node {
	return h.Div(
		h.Class(styles.ErrorBanner),
		h.P(g.Text("Error: "+msg)),
	)
}
