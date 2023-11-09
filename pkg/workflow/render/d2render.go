package render

import (
	"context"
	"os"
	"oss.terrastruct.com/d2/d2compiler"
	"oss.terrastruct.com/d2/d2exporter"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"path/filepath"
	"stories/pkg/workflow/common"
	"strings"
)

const graphDefinition = `
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
`

// Remember to add if err != nil checks in production.
func RenderWorkflow(wf common.Workflow) {
	graph, config, _ := d2compiler.Compile("", strings.NewReader(graphDefinition), nil)
	graph.ApplyTheme(d2themescatalog.NeutralDefault.ID)
	ruler, _ := textmeasure.NewRuler()
	_ = graph.SetDimensions(nil, ruler, nil)
	_ = d2dagrelayout.Layout(context.Background(), graph, nil)
	diagram, _ := d2exporter.Export(context.Background(), graph, nil)
	diagram.Config = config
	out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		ThemeID: &d2themescatalog.NeutralDefault.ID,
	})
	_ = os.WriteFile(filepath.Join("out.svg"), out, 0600)
}
