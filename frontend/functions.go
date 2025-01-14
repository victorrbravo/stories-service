package frontend

type ComponentInterface interface {
	Type() uint8
}

const (
	Red uint8 = iota + 1
	Blue
	Brown
	Black
)

const (
	TypeP uint8 = iota + 1
	TypeDiv
	TypeButton
	TypeH1
	TypeH2
	TypeH3
	TypeH4
	TypeH5
	TypeH6
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

type H2Properties PProperties

type H2Struct struct {
	Content    string
	Properties H2Properties
}

func H2(content string, properties H2Properties) *H2Struct {
	return &H2Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H2Struct) Type() uint8 { return TypeH2 }

type H3Properties PProperties

type H3Struct struct {
	Content    string
	Properties H3Properties
}

func H3(content string, properties H3Properties) *H3Struct {
	return &H3Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H3Struct) Type() uint8 { return TypeH3 }

type H4Properties PProperties

type H4Struct struct {
	Content    string
	Properties H4Properties
}

func H4(content string, properties H4Properties) *H4Struct {
	return &H4Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H4Struct) Type() uint8 { return TypeH4 }

type H5Properties PProperties

type H5Struct struct {
	Content    string
	Properties H5Properties
}

func H5(content string, properties H5Properties) *H5Struct {
	return &H5Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H5Struct) Type() uint8 { return TypeH5 }

type H6Properties PProperties

type H6Struct struct {
	Content    string
	Properties H6Properties
}

func H6(content string, properties H6Properties) *H6Struct {
	return &H6Struct{
		Content:    content,
		Properties: properties,
	}
}

func (h *H6Struct) Type() uint8 { return TypeH6 }

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

type DivStruct struct {
	Components []ComponentInterface
}

func (d *DivStruct) Type() uint8 {
	return TypeDiv
}

func Div(components ...ComponentInterface) *DivStruct {
	return &DivStruct{
		Components: components,
	}
}
