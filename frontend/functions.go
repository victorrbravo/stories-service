package frontend

type ComponentInterface interface {
	Type() uint8
}

const (
	Red uint8 = iota + 1
	Blue
	Brown
	Dark
)

const (
	TypeText uint8 = iota + 1
	TypeView
	TypeButton
)

type TextProperties struct {
	Color uint8
}

type TextStruct struct {
	Content    string
	Properties TextProperties
}

func Text(content string, properties TextProperties) *TextStruct {
	return &TextStruct{
		Content:    content,
		Properties: properties,
	}
}

func (t *TextStruct) Type() uint8 {
	return TypeText
}

type ButtonProperties TextProperties

type ButtonStruct struct {
	Text       string
	Properties ButtonProperties
}

func Button(text string, properties ButtonProperties) *ButtonStruct {
	return &ButtonStruct{
		Text:       text,
		Properties: properties,
	}
}

func (b *ButtonStruct) Type() uint8 {
	return TypeButton
}

type ViewStruct struct {
	Components []ComponentInterface
}

func (v *ViewStruct) Type() uint8 {
	return TypeView
}

func View(components ...ComponentInterface) *ViewStruct {
	return &ViewStruct{
		Components: components,
	}
}
