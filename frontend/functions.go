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
	TypeP uint8 = iota + 1
	TypeView
	TypeButton
	TypeH1
)

type PProperties struct {
	Color uint8
}

type PStruct struct {
	Content    string
	Properties PProperties
}

func P(content string, properties PProperties) *PStruct {
	return &PStruct{
		Content:    content,
		Properties: properties,
	}
}

func (t *PStruct) Type() uint8 {
	return TypeP
}

type H1Properties PProperties

type H1Struct struct {
	Content    string
	Properties H1Properties
}

func H1(content string, properties H1Properties) *H1Struct {
	return &H1Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H1Struct) Type() uint8 {
	return TypeH1
}

type ButtonProperties PProperties

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
