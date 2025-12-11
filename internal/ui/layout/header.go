package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Header(title, description string) g.Node {
	return h.Header(
		h.Class("header"),
		h.H1(h.Class("header__title"), g.Text(title)),
		h.P(h.Class("header__description"), g.Text(description)),
	)
}
