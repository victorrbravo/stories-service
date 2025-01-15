package main

import (
	f "stories/frontend"
)

func DoubleRedText(content string, secondContent string) f.ComponentInterface {
	return f.Div(
		f.P(content, f.PProperties{
			Color: f.Red,
		}),
		f.P(secondContent, f.PProperties{
			Color: f.Red,
		}),
	)
}

func BrownButton(content string) f.ComponentInterface {
	return f.Button(content, f.ButtonProperties{
		Color: f.Brown,
	})
}

func DoubleBlueButtonsWithEqualText(content string) f.ComponentInterface {
	return f.Div(
		f.Button(content, f.ButtonProperties{
			Color: f.Blue,
		}),
		f.Button(content, f.ButtonProperties{
			Color: f.Blue,
		}),
	)
}

func ReturnView() *f.DivStruct {
	return f.Div(
		f.H1("Título de la página", f.H1Properties{
			Color: f.Red,
		}),
		f.H2("Subtítulo de la página", f.H2Properties{
			Color: f.Black,
			Class: "mi h2cito",
		}),
		f.H3("Ejemplo de H3", f.H3Properties{
			Color: f.Brown,
		}),
		f.H4("Ejemplo de H4", f.H4Properties{
			Color: f.Blue,
		}),
		f.H5("H5 sin color predefinido", f.H5Properties{}),
		f.H6("H6 Marrón", f.H6Properties{
			Color: f.Brown,
		}),
		f.P("Hello World!", f.PProperties{
			Color: f.Blue,
		}),
		DoubleRedText("Hola Mundo!", "こんいちわ！"),
		f.Button("Presiona aquí", f.ButtonProperties{
			Color: f.Red,
		}),
		BrownButton("Presiona aquí de nuevo"),
		DoubleBlueButtonsWithEqualText("Somos dos botones azules"),
	)
}

func main() {
	view := ReturnView()
	f.GenerateHTMLCode(f.PageObject{
		PageName:     "Example",
		PageLanguage: "es",
		Components:   *view,
	})
}
