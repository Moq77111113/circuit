package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/ui/assets"
)

func renderHead(title string) []g.Node {
	return []g.Node{
		h.Head(
			h.Meta(h.Charset("utf-8")),
			h.Meta(
				h.Name("viewport"),
				h.Content("width=device-width, initial-scale=1"),
			),
			h.Link(h.Rel("icon"), h.Type("image/svg+xml"), h.Href("data:image/svg+xml,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%20100%20100'%3E%3Cdefs%3E%3Cfilter%20id='glow'%3E%3CfeGaussianBlur%20stdDeviation='2'%20result='blur'/%3E%3CfeMerge%3E%3CfeMergeNode%20in='blur'/%3E%3CfeMergeNode%20in='SourceGraphic'/%3E%3C/feMerge%3E%3C/filter%3E%3C/defs%3E%3Crect%20width='100'%20height='100'%20fill='%23000000'/%3E%3Cpath%20d='M0%2050%20L10%2030%20L20%2055%20L30%2020%20L40%2060%20L50%2025%20L60%2055%20L70%2015%20L80%2050%20L90%2035%20L100%2050'%20stroke='%23FFFFFF'%20stroke-width='6'%20stroke-linecap='round'%20fill='none'%20filter='url(%23glow)'/%3E%3C/svg%3E")),
			h.TitleEl(g.Text(title)),
			h.StyleEl(g.Raw(assets.DefaultCSS)),
		),
	}
}
