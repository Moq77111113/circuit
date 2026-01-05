package tags

type Option struct {
	Label string
	Value string
}

type Field struct {
	Name        string
	Type        string
	InputType   InputType
	Help        string
	Required    bool
	ReadOnly    bool
	Min         string
	Max         string
	Step        string
	Pattern     string
	MinLen      int
	MaxLen      int
	Options     []Option
	Fields      []Field
	IsSlice     bool
	ElementType string
}
