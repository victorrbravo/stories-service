package pg

import (
	"context"
	"fmt"
	"log"
	"os"
	"oss.terrastruct.com/d2/d2compiler"
	"oss.terrastruct.com/d2/d2exporter"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"path/filepath"
	db "stories/database"
	"stories/pkg/workflow/common"
	"strings"
)

func GenerateGraphDefinition[T db.Scanner](wf *common.Workflow[T]) (string, error) {
	var str string
	var styles string
	ctx := context.Background()
	statuses, err := wf.Statuses(ctx, false)
	dataForTasks := ""

	if err != nil {
		return str, err
	}

	for _, name := range statuses {
		current, ok := wf.Tasks[name]
		if !ok {
			log.Println("Error accessing name", name)
			continue
		}
		str = str + name
		str = str + common.ARROW_TASK
		newStyle := fmt.Sprintf(`
		%s.style.border-radius: %s
		%s.style.fill: "%s"
`, name, "8", name, SelectColor(current.Data.Count, wf.Data.Stats.Min, wf.Data.Stats.Max))
		styles = styles + newStyle
		titleTime := current.Data.TimeField.String()
		titleRole := current.Data.InfoField.String()
		titleCount := current.Data.Count.String()
		title := titleTime + "\\n" + titleRole + "\\n" + titleCount
		dataForTasks += fmt.Sprintf(`
"%s" -> %s: {
style: {
	stroke-dash: 3
}
}
"%s".style.stroke: gray
"%s".style.fill: white
"%s".style.font-size: 12
"%s".style.font-color: gray
"%s".style.shadow: true

`, title, name, title, title, title, title, title)

	}
	str, _ = strings.CutSuffix(str, common.ARROW_TASK)
	str += "\n"
	str += styles
	str += "\n"
	str += dataForTasks
	return str, nil
}

func DummyFunc[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		return wf, nil
	}
}

func Render[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	definition, err := GenerateGraphDefinition(wf)
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		if err != nil {
			return nil, err
		}
		graph, config, _ := d2compiler.Compile("", strings.NewReader(definition), nil)
		graph.ApplyTheme(d2themescatalog.NeutralDefault.ID)
		ruler, _ := textmeasure.NewRuler()
		_ = graph.SetDimensions(nil, ruler, nil)
		_ = d2dagrelayout.Layout(context.Background(), graph, nil)
		diagram, _ := d2exporter.Export(context.Background(), graph, nil)
		diagram.Config = config
		out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
			ThemeID: &d2themescatalog.NeutralDefault.ID,
		})
		_ = os.WriteFile(filepath.Join("workflow.svg"), out, 0600)
		fmt.Println("workflow written")
		return wf, nil
	}
}
