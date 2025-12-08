package tags

type InputType string

const (
	TypeText     InputType = "text"
	TypePassword InputType = "password"
	TypeNumber   InputType = "number"
	TypeCheckbox InputType = "checkbox"
	TypeRadio    InputType = "radio"
	TypeRange    InputType = "range"
	TypeDate     InputType = "date"
	TypeTime     InputType = "time"
	TypeEmail    InputType = "email"
	TypeTel      InputType = "tel"
	TypeUrl      InputType = "url"
	TypeColor    InputType = "color"
	TypeFile     InputType = "file"
	TypeHidden   InputType = "hidden"
	TypeSelect   InputType = "select"
	TypeSection  InputType = "section"
)
