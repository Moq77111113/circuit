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

func Footer() g.Node {
	return h.Footer(
		h.Class("footer"),
		h.P(
			g.Text("Powered by "),
			h.A(
				h.Href("https://github.com/moq77111113/circuit"),
				h.Target("_blank"),
				h.Rel("noopener noreferrer"),
				h.Class("footer__link"),
				g.Text("Circuit"),
			),
		),
	)
}

func FavIcons() g.Node {

	return g.Group([]g.Node{
		h.Link(h.Rel("icon"), h.Type("image/svg+xml"), h.Href("data:image/svg+xml,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%20100%20100'%3E%3Cdefs%3E%3Cfilter%20id='glow'%3E%3CfeGaussianBlur%20stdDeviation='2'%20result='blur'/%3E%3CfeMerge%3E%3CfeMergeNode%20in='blur'/%3E%3CfeMergeNode%20in='SourceGraphic'/%3E%3C/feMerge%3E%3C/filter%3E%3C/defs%3E%3Crect%20width='100'%20height='100'%20fill='%23000000'/%3E%3Cpath%20d='M0%2050%20L10%2030%20L20%2055%20L30%2020%20L40%2060%20L50%2025%20L60%2055%20L70%2015%20L80%2050%20L90%2035%20L100%2050'%20stroke='%23FFFFFF'%20stroke-width='6'%20stroke-linecap='round'%20fill='none'%20filter='url(%23glow)'/%3E%3C/svg%3E")),
	})
}
