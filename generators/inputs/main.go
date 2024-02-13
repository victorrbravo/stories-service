package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"stories/pkg/workflow/common"
	"strings"
	"text/template"

	"github.com/twinj/uuid"

	"stories/generators"
)

type CustomString string

type entity struct {
	Model        CustomString
	TableName    CustomString
	Schema       CustomString
	UsesEnums    bool
	UsesTypes    bool
	NeedsPostgis bool
	Attrs        []attr
	Searcher     []attr
	OrderBy      CustomString
	Order        CustomString
	AddPaging    bool
	Models       []Description
}

type Description struct {
	Name        CustomString
	KeyField    CustomString
	Searchables []attr
}

type attr struct {
	Name       CustomString
	FieldName  CustomString
	Key        bool
	Type       string
	Nullable   bool
	Scannable  bool
	Editable   bool
	Searchable bool
	Discrete   bool
	KeyForMany bool
	Getter     bool
	Title      CustomString
}

type generatorContext struct {
	template   *template.Template
	functions  *template.Template
	attrs      []attr
	searcher   []attr
	noTypes    bool
	model      string
	tableName  string
	schema     string
	orderBy    string
	order      string
	addPaging  bool
	sourcePath string
	targetPath string
	models     []Description
	keys       map[string]string
}

const (
	ProcessedFillingData    = "OrgFilingData"
	ProcessedIndFillingData = "IndFilingData"
	ProcessedCtyFillingData = "CtyFilingData"
	UserWhitelist           = "UserWhitelist"
	News                    = "News"
	SuffixData              = "Data"
	Analysis                = "SwotAnalysis"
	charset                 = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func randomInteger() int {
	return rand.Intn(100)
}

func randomString() string {
	var data []byte = make([]byte, 10)
	for index := range data {
		data[index] += charset[rand.Intn(len(charset))]
	}
	return string(data)
}

func randomArrayOfStrings() []string {
	return []string{
		randomString(),
		randomString(),
	}
}

func (a attr) Random() string {
	switch a.Type {
	case "types.String":
		return fmt.Sprintf("types.NewString(\"%s\")", randomString())
	case "types.UUID":
		return fmt.Sprintf("types.NewUUIDMust(\"%s\")", uuid.NewV4().String())
	case "types.Int":
		return fmt.Sprintf("types.NewInt(%d)", randomInteger())
	case "types.DateTime":
		return "types.DateTime(time.Now())"
	case "types.StringArray":
		return fmt.Sprintf("[]string{\"%s\"}", strings.Join(randomArrayOfStrings(), "\", \""))
	case "int":
		return fmt.Sprintf("%d", randomInteger())
	case "string":
		return fmt.Sprintf("\"%s\"", randomString())
	case "types.PortfolioType":
		// FIXME: should be able to acknowledge enums
		return "\"Public\""
	case "float32":
		return fmt.Sprintf("%.8f", rand.Float32()+0.5)
	case "postgis.Point":
		return ""
	default:
		// panic(fmt.Errorf("cannot generate random sample for type `%s'", a.Type))
		return ""
	}
}

func (e entity) KeyType() string {
	for _, attr := range e.Attrs {
		if attr.Key {
			return strings.TrimPrefix(attr.Type, "types.")
		}
	}
	panic(fmt.Errorf("cannot find a key for `%s'", e.Model))
}

func shouldImportEnums(attrs []attr) bool {
	for _, each := range attrs {
		if strings.HasPrefix(each.Type, "enums.") {
			return true
		}
	}
	return false
}

func shouldImportPostgis(attrs []attr) bool {
	for _, each := range attrs {
		if strings.HasPrefix(each.Type, "postgis.") {
			return true
		}
	}
	return false
}

func shouldImportTypes(attrs []attr) bool {
	for _, each := range attrs {
		if strings.HasPrefix(each.Type, "types.") {
			return true
		}
	}
	return false
}

func readAllTextFromFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		innerErr := file.Close()
		if innerErr != nil {
			log.Println(innerErr)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func onlyNonKeys(attrs []attr) []attr {
	var array []attr = make([]attr, 0)
	for _, each := range attrs {
		if each.Key && each.Editable == false {
			continue
		}
		array = append(array, each)
	}
	return array
}

func checkIfHasFieldNameForUpdateOrDelete(fieldname CustomString, name CustomString) CustomString {
	var title string = fmt.Sprintf(`\"%s\"`, string(name.ToSnakeCase()))
	var fieldTitle string = strings.Trim(string(fieldname.ToSnakeCase()), " ")
	if len(fieldTitle) > 0 {
		title = fmt.Sprintf(`\"%s\"`, fieldTitle)
	}
	return CustomString(title)
}

func checkIfHasFieldNameForSelect(fieldname CustomString, name CustomString) CustomString {
	var title string = fmt.Sprintf(`\"%s\"`, string(name.ToSnakeCase()))
	var fieldTitle string = strings.Trim(string(fieldname.ToSnakeCase()), " ")
	if len(fieldTitle) > 0 {
		title = fmt.Sprintf(`\"%s\" AS \"%s\"`, fieldTitle, string(name.ToSnakeCase()))
	}
	return CustomString(title)
}

func onlyKeys(attrs []attr) []attr {
	var array []attr = make([]attr, 0)
	for _, each := range attrs {
		if !each.Key {
			continue
		}
		array = append(array, each)
	}
	return array
}

func (str CustomString) Plural() string {
	var lastIndex int = len(str) - 1
	switch str {
	case ProcessedIndFillingData, ProcessedCtyFillingData, Analysis, News, UserWhitelist:
		return string(str)
	case ProcessedFillingData:
		return string(str)
	default:
		if strings.HasSuffix(string(str), SuffixData) {
			return string(str)
		}
	}
	if str[lastIndex] == 'y' {
		return string(str)[:lastIndex] + "ies"
	} else {
		return string(str) + "s"
	}
}

func (str CustomString) ToYamlTag() string {
	return fmt.Sprintf("`yaml:\"%s\"`", str.ToSnakeCase())
}

func (str CustomString) ToJsonTag() string {
	return fmt.Sprintf("`json:\"%s\"`", str.ToSnakeCase())
}

func (str CustomString) ToCamelCase() string {
	re := regexp.MustCompile("([A-Z][a-z]+|ID)")
	var matches [][]string = re.FindAllStringSubmatch(string(str), -1)
	if len(matches) == 0 {
		return strings.ToLower(string(str))
	}
	var parts []string = make([]string, len(matches))
	for index, match := range matches {
		parts[index] = match[1]
	}
	parts[0] = strings.ToLower(parts[0])
	// Return the joined thing
	camelCase := strings.Join(parts, "")
	switch camelCase {
	case "type": // FIXME: add any other keywords here
		return "itemType"
	default:
		return camelCase
	}
}

func (str CustomString) ToSnakeCases() string {
	value := string(str.Plural())
	re := regexp.MustCompile("([A-Z][a-z]+|ID)")
	var matches [][]string = re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return strings.ToLower(value)
	}
	var parts []string = make([]string, len(matches))
	for index, match := range matches {
		parts[index] = strings.ToLower(match[1])
	}
	// Return the joined thing
	return strings.Join(parts, "_")
}

func (str CustomString) ToSnakeCase() string {
	re := regexp.MustCompile("([A-Z][a-z]+|ID)")
	var matches [][]string = re.FindAllStringSubmatch(string(str), -1)
	if len(matches) == 0 {
		return strings.ToLower(string(str))
	}
	var parts []string = make([]string, len(matches))
	for index, match := range matches {
		parts[index] = strings.ToLower(match[1])
	}
	// Return the joined thing
	return strings.Join(parts, "_")
}

func (str CustomString) Simplified() string {
	value := str.ToSnakeCase()
	value, _ = strings.CutPrefix(value, common.TablePrefix)
	return value
}

func checkStringIsNotEmpty(str CustomString) bool {
	if len(string(str)) > 0 {
		return true
	}
	return false
}

func (context *generatorContext) SetSourcePath(sourcePath string) {
	context.sourcePath = sourcePath
}

func (context *generatorContext) GetSourcePath() string {
	return context.sourcePath
}

func (context *generatorContext) GetTargetPath() string {
	return context.targetPath
}

func (context *generatorContext) Prepare() {
	parser1 := template.New("inputs")
	// Initialize all
	funcMap := template.FuncMap{
		"onlyNonKeys": onlyNonKeys,
		"add": func(numbers ...int) int {
			var total int = 0
			for _, num := range numbers {
				total += num
			}
			return total
		},
		"checkForSelect":      checkIfHasFieldNameForSelect,
		"checkUpdateOrDelete": checkIfHasFieldNameForUpdateOrDelete,
		"onlyKeys":            onlyKeys,
		"checkIsNotEmpty":     checkStringIsNotEmpty,
	}
	parser2 := template.New("functions")
	tplContent, err := readAllTextFromFile("../../generators/inputs/template.tpl")
	if err != nil {
		return
	}

	funcContent, err := readAllTextFromFile("../../generators/inputs/functions.tpl")
	if err != nil {
		return
	}
	context.template = template.Must(parser1.Funcs(funcMap).Parse(string(tplContent)))
	context.functions = template.Must(parser2.Funcs(funcMap).Parse(string(funcContent)))

}

func (context *generatorContext) testAlreadyExists(fileName string) bool {
	var fullName string
	var filePath string
	fullName = fmt.Sprintf("%s_test.go", fileName)
	filePath = path.Join(context.GetTargetPath(), path.Join("tests", fullName))
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else {
		return false
	}
}

func (context *generatorContext) Write() error {
	var file *os.File
	var name string
	var err error
	var fileName string
	fileName = generators.GetFileName(context.model)
	name = fmt.Sprintf("%s.go", fileName)
	file, err = os.Create(path.Join(context.GetTargetPath(), name))
	if err != nil {
		return err
	}
	if len(context.tableName) == 0 {
		context.tableName = context.model
	}
	err = context.template.Execute(file, entity{
		Model:        CustomString(context.model),
		TableName:    CustomString(context.tableName),
		Schema:       CustomString(context.schema),
		UsesTypes:    shouldImportTypes(context.attrs),
		UsesEnums:    shouldImportEnums(context.attrs),
		NeedsPostgis: shouldImportPostgis(context.attrs),
		Attrs:        context.attrs,
		Searcher:     context.searcher,
		OrderBy:      CustomString(context.orderBy),
		Order:        CustomString(context.order),
		AddPaging:    context.addPaging,
	})
	context.tableName = ""
	context.orderBy = ""
	context.order = ""
	context.addPaging = false
	if err != nil {
		return err
	}
	return nil
}

func (context *generatorContext) WriteFunctionsFile() error {
	var file *os.File
	var name string
	var err error
	var fileName string
	fileName = "functions"
	name = fmt.Sprintf("%s.go", fileName)
	file, err = os.Create(path.Join(context.GetTargetPath(), name))
	if err != nil {
		return err
	}
	if len(context.tableName) == 0 {
		context.tableName = context.model
	}
	err = context.functions.Execute(file, entity{
		Model:     CustomString(context.model),
		TableName: CustomString(context.tableName),
		Schema:    CustomString(context.schema),
		UsesTypes: shouldImportTypes(context.attrs),
		UsesEnums: shouldImportEnums(context.attrs),
		Attrs:     context.attrs,
		Searcher:  context.searcher,
		Models:    context.models,
	})
	if err != nil {
		return err
	}
	return nil
}

func (context *generatorContext) HandleFile(file os.FileInfo) error {
	var entity entity
	context.attrs = make([]attr, 0)
	context.searcher = make([]attr, 0)
	if context.models == nil {
		context.models = make([]Description, 0)
	}
	if context.keys == nil {
		context.keys = make(map[string]string)
	}
	filePath := context.sourcePath + "/" + file.Name()
	// Read all the content of the file
	content, err := readAllTextFromFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &entity)
	if err != nil {
		return err
	}
	context.model = string(entity.Model)
	if len(string(entity.Schema)) > 0 {
		context.schema = fmt.Sprintf("%s.", strings.TrimRight(string(entity.Schema), "."))
	} else {
		context.schema = ""
	}
	if len(string(entity.TableName)) > 0 {
		context.tableName = string(entity.TableName)
	} else {
		context.tableName = context.model
	}
	context.orderBy = string(entity.OrderBy)
	context.order = string(entity.Order)
	context.addPaging = entity.AddPaging
	searchables := []attr{}
	for _, curAttr := range entity.Attrs {
		if curAttr.Searchable {
			searchables = append(searchables, curAttr)
		}
	}
	for _, curAttr := range entity.Attrs {
		context.attrs = append(context.attrs, curAttr)
		if curAttr.Key {
			desc := Description{entity.Model, curAttr.Name, searchables}
			context.models = append(context.models, desc)
		}
	}
	for _, attr := range entity.Searcher {
		context.searcher = append(context.searcher, attr)
	}
	return nil
}

func main() {
	targetPath := "../../generated/inputs"
	var files []os.FileInfo
	var err error
	handler := &generatorContext{targetPath: targetPath}
	// Let the handler prepare its
	// elf
	handler.Prepare()
	handler.SetSourcePath("../../specs")
	files, err = generators.GetSourceFiles(handler)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files {
		// Generate
		err = generators.GenerateFile(file, handler)
		if err != nil {
			log.Println(err)
		}
		err = handler.Write()
		if err != nil {
			log.Fatalln(err)
		}
	}
	handler.WriteFunctionsFile()
}
