package main

import (
	f "stories/frontend"
)

func DoubleRedText(content string, secondContent string) f.ComponentInterface {
	return f.View(
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
	return f.View(
		f.Button(content, f.ButtonProperties{
			Color: f.Blue,
		}),
		f.Button(content, f.ButtonProperties{
			Color: f.Blue,
		}),
	)
}

func ReturnView() *f.ViewStruct {
	return f.View(
		f.H1("Título de la página", f.H1Properties{
			Color: f.Red,
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
