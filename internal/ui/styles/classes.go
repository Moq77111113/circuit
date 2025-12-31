package styles

const (
	// Collapsible component
	Collapsible          = "collapsible"
	CollapsibleHeader    = "collapsible__header"
	CollapsibleTitle     = "collapsible__title"
	CollapsibleIcon      = "collapsible__icon"
	CollapsibleCount     = "collapsible__count"
	CollapsibleSummary   = "collapsible__summary"
	CollapsibleBody      = "collapsible__body"
	CollapsibleCollapsed = "collapsible--collapsed"
	CollapsibleDepth0    = "collapsible--depth-0"
	CollapsibleDepth1    = "collapsible--depth-1"
	CollapsibleDepth2    = "collapsible--depth-2"
	CollapsibleDepth3    = "collapsible--depth-3"
	CollapsibleDepth4    = "collapsible--depth-4"

	// Field component
	Field           = "field"
	FieldLabel      = "field__label"
	FieldInput      = "field__input"
	FieldHelp       = "field__help"
	FieldError      = "field__error"
	FieldRequired   = "field__label-required"
	FieldLabelClick = "field__label--clickable"
	FieldSelect     = "field__select"

	// Button component
	Button          = "button"
	ButtonPrimary   = "button--primary"
	ButtonSecondary = "button--secondary"
	ButtonDanger    = "button--danger"
	ButtonAdd       = "button--add"
	ButtonRemove    = "button--remove"

	// Card component
	Card        = "card"
	CardHeader  = "card__header"
	CardTitle   = "card__title"
	CardBody    = "card__body"
	CardIcon    = "card__icon"
	CardName    = "card__name"
	CardArrow   = "card__arrow"
	CardPreview = "card__preview"

	// Struct card (specific type of card)
	StructCard        = "struct-card"
	StructCardHeader  = "struct-card__header"
	StructCardName    = "struct-card__name"
	StructCardArrow   = "struct-card__arrow"
	StructCardPreview = "struct-card__preview"

	// Nav component
	Nav                = "nav"
	NavTitle           = "nav__title"
	NavList            = "nav__list"
	NavSublist         = "nav__sublist"
	NavItem            = "nav__item"
	NavItemCollapsible = "nav__item--collapsible"
	NavItemNested      = "nav__item--nested"
	NavLink            = "nav__link"
	NavLinkField       = "nav__link--field"
	NavLinkSection     = "nav__link--section"
	NavLinkSlice       = "nav__link--slice"
	NavLinkSliceItem   = "nav__link--slice-item"
	NavChevron         = "nav__chevron"

	// Icon modifiers
	Icon             = "icon"
	IconChevronDown  = "icon--chevron-down"
	IconChevronRight = "icon--chevron-right"
	IconArrowRight   = "icon--arrow-right"

	// Slice components
	SliceItem             = "slice__item"
	SliceItemPrimitive    = "slice__item--primitive"
	SliceItemStruct       = "slice__item--struct"
	SliceItemHeader       = "slice__item-header"
	SliceItemBody         = "slice__item-body"
	SliceItemTitle        = "slice__item-title"
	SliceItemRemoveButton = "slice-item__remove-button"
	SliceChevron          = "slice__chevron"
	SliceAddButton        = "slice__add-button"

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
	Breadcrumb          = "breadcrumb"
	BreadcrumbItem      = "breadcrumb__item"
	BreadcrumbLink      = "breadcrumb__link"
	BreadcrumbLinkRoot  = "breadcrumb__link--root"
	BreadcrumbSeparator = "breadcrumb__separator"
	BreadcrumbIndex     = "breadcrumb__index"

	// Header
	Header            = "header"
	HeaderTitle       = "header__title"
	HeaderDescription = "header__description"

	// Footer
	Footer     = "footer"
	FooterLink = "footer__link"

	// Mobile
	MobileMenuToggle = "mobile-menu-toggle"
	MobileOverlay    = "mobile-overlay"

	// Input components
	RangeWrapper = "range-wrapper"
	RangeMin     = "range-min"
	RangeMax     = "range-max"
	RangeValue   = "range-value"

	RadioGroup  = "radio-group"
	RadioOption = "radio-option"

	ToggleSwitch         = "toggle-switch"
	ToggleSwitchInput    = "toggle-switch__input"
	ToggleSwitchLabel    = "toggle-switch__label"
	ToggleSwitchLabelOn  = "toggle-switch__label--on"
	ToggleSwitchLabelOff = "toggle-switch__label--off"
	ToggleSwitchSlider   = "toggle-switch__slider"

	// Sidebar tree
	SidebarTree  = "sidebar-tree"
	TreeRoot     = "tree-root"
	TreeNode     = "tree-node"
	TreeNodeLink = "tree-node__link"

	// Actions dropdown
	ActionsDropdown      = "actions-dropdown"
	ActionsButton        = "actions-button"
	ActionsIcon          = "actions-icon"
	ActionsMenu          = "actions-menu"
	ActionsMenuItem      = "actions-menu__item"
	ActionsMenuItemForm  = "actions-menu__item-form"
	ActionsMenuItemLabel = "actions-menu__item-label"
	ActionsMenuItemDesc  = "actions-menu__item-desc"

	// Error banner
	ErrorBanner = "error-banner"

	// State and misc
	EmptyState = "empty-state"
	Collapsed  = "collapsed"
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

// Merge combines multiple CSS class names into a single space-separated string.
// Empty strings are ignored.
func Merge(classes ...string) string {
	if len(classes) == 0 {
		return ""
	}
	if len(classes) == 1 {
		return classes[0]
	}

	// Pre-calculate total length to avoid multiple allocations
	totalLen := 0
	nonEmpty := 0
	for _, c := range classes {
		if c != "" {
			totalLen += len(c)
			nonEmpty++
		}
	}
	if nonEmpty == 0 {
		return ""
	}

	// Add space separators
	totalLen += nonEmpty - 1

	// Build string efficiently
	result := make([]byte, 0, totalLen)
	first := true
	for _, c := range classes {
		if c != "" {
			if !first {
				result = append(result, ' ')
			}
			result = append(result, c...)
			first = false
		}
	}
	return string(result)
}
