package main

import (
	f "stories/frontend"
)

func DoubleRedText(content string, secondContent string) f.ComponentInterface {
	return f.View(
		f.Text(content, f.TextProperties{
			Color: f.Red,
		}),
		f.Text(secondContent, f.TextProperties{
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
		f.Text("Hello World!", f.TextProperties{
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

// func RecordComponents(components f.ViewStruct) {
// 	for i := range components.Components {
// 		switch component := components.Components[i].(type) {
// 		case *f.TextStruct:
// 			fmt.Println(component.Content, component.Properties.Color, component.Type())
// 		case *f.ButtonStruct:
// 			component.Properties.OnPress()
// 			fmt.Println(component.Text, component.Properties.Color, component.Type())
// 		case *f.ViewStruct:
// 			RecordComponents(*component)
// 		}
// 	}
// }

func main() {
	view := ReturnView()
	f.GenerateHTMLCode(f.PageObject{
		PageName:     "Example",
		PageLanguage: "es",
		Components:   *view,
	})
}
