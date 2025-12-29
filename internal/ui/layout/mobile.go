package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func renderMobileToggle() g.Node {
	return h.Button(
		h.Class("mobile-menu-toggle"),
		h.Type("button"),
		g.Attr("aria-label", "Toggle menu"),
		g.Attr("onclick", "toggleSidebar()"),
		g.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M9 6l6 6-6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`),
	)
}

func renderMobileOverlay() g.Node {
	return h.Div(
		h.Class("mobile-overlay"),
		g.Attr("onclick", "closeSidebar()"),
	)
}
