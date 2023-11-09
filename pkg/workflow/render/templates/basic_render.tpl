{{- range .Tasks }}
	{{ .Name }} {{ .Type }} {{ .Name.ToJsonTag }}
{{- end }}

ToDo -> Progress -> Done
tasks -> Progress:  {
style: {
	stroke-dash: 3
}
}
"role: one\n3min" -> ToDo:  {
style: {
	stroke-dash: 3
}
}
tasks3-> Done:  {
style: {
	stroke-dash: 3
}
}
tasks.style.stroke: "#f4a261"
tasks.style.fill: white
"role: one\n3min".style.stroke: "#f4a261"
"role: one\n3min".style.fill: white
tasks3.style.stroke: "#f4a261"
tasks3.style.fill: white

tasks.style.font-size: 10
"role: one\n3min".style.font-size: 10
"role: one\n3min".style.font-color: gray
tasks3.style.font-size: 10

ToDo.style.border-radius: 8
Progress.style.border-radius: 8
Done.style.border-radius: 8