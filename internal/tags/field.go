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
	Min         string
	Max         string
	Step        string
	Options     []Option
	Fields      []Field
	IsSlice     bool
	ElementType string
}
