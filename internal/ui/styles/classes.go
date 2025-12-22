package styles

// BEM-style CSS class constants
// Block__Element--Modifier naming convention

const (
	// Collapsible component
	Collapsible          = "collapsible"
	CollapsibleHeader    = "collapsible__header"
	CollapsibleTitle     = "collapsible__title"
	CollapsibleIcon      = "collapsible__icon"
	CollapsibleCount     = "collapsible__count"
	CollapsibleBody      = "collapsible__body"
	CollapsibleCollapsed = "collapsible--collapsed"
	CollapsibleDepth0    = "collapsible--depth-0"
	CollapsibleDepth1    = "collapsible--depth-1"
	CollapsibleDepth2    = "collapsible--depth-2"
	CollapsibleDepth3    = "collapsible--depth-3"
	CollapsibleDepth4    = "collapsible--depth-4"

	// Field component
	Field            = "field"
	FieldLabel       = "field__label"
	FieldInput       = "field__input"
	FieldHelp        = "field__help"
	FieldRequired    = "field__label-required"
	FieldLabelClick  = "field__label--clickable"
	FieldSelect      = "field__select"

	// Button component
	Button        = "button"
	ButtonPrimary = "button--primary"

	// Card component
	Card       = "card"
	CardHeader = "card__header"
	CardTitle  = "card__title"
	CardBody   = "card__body"
	CardIcon   = "card__icon"

	// Nav component
	Nav        = "nav"
	NavTitle   = "nav__title"
	NavList    = "nav__list"
	NavSublist = "nav__sublist"
	NavItem    = "nav__item"
	NavLink    = "nav__link"
	NavChevron = "nav__chevron"

	// Icon modifiers
	Icon             = "icon"
	IconChevronDown  = "icon--chevron-down"
	IconChevronRight = "icon--chevron-right"
	IconArrowRight   = "icon--arrow-right"

	// Form
	Form        = "form"
	FormSection = "form__section"
	FormActions = "form__actions"

	// Layout
	App          = "app"
	AppSidebar   = "app__sidebar"
	AppMain      = "app__main"
	AppContainer = "app__container"

	// Breadcrumb
	Breadcrumb     = "breadcrumb"
	BreadcrumbItem = "breadcrumb__item"
	BreadcrumbLink = "breadcrumb__link"

	// State modifiers
	Collapsed = "collapsed"
)

// DepthClass returns the BEM depth modifier class for a given depth level.
// Clamps depth to range [0, 4].
func DepthClass(depth int) string {
	if depth < 0 {
		depth = 0
	}
	if depth > 4 {
		depth = 4
	}
	return []string{
		CollapsibleDepth0,
		CollapsibleDepth1,
		CollapsibleDepth2,
		CollapsibleDepth3,
		CollapsibleDepth4,
	}[depth]
}
