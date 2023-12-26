// nolint: govet
package pg

import (
	"context"
	"errors"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"log"
	"reflect"
	db "stories/database"
	"stories/pkg/workflow/common"
	"strings"
)

const (
	AndOperator = " AND "
	EndTask     = "END"
)

type Boolean bool

type WhereTask struct {
	Name  string
	Where string
}

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Top        *Term             `"SELECT" ( "TOP" @@ )?`
	Distinct   bool              `(  @"DISTINCT"`
	All        bool              ` | @"ALL" )?`
	Expression *SelectExpression `@@`
	From       *From             `"FROM" @@`
	Limit      *Expression       `( "LIMIT" @@ )?`
	Offset     *Expression       `( "OFFSET" @@ )?`
	GroupBy    *Expression       `( "GROUP" "BY" @@ )?`
}

type From struct {
	TableExpressions []*TableExpression `@@ ( "," @@ )*`
	Where            *Expression        `( "WHERE" @@ )?`
}

type TableExpression struct {
	Table  string        `( @Ident ( "." @Ident )*`
	Select *Select       `  | "(" @@ ")"`
	Values []*Expression `  | "VALUES" "(" @@ ( "," @@ )* ")")`
	As     string        `( "AS" @Ident )?`
}

type SelectExpression struct {
	All         bool                 `  @"*"`
	Expressions []*AliasedExpression `| @@ ( "," @@ )*`
}

type AliasedExpression struct {
	Expression *Expression `@@`
	As         string      `( "AS" @Ident )?`
}

type Expression struct {
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

func (o OrCondition) GetAnd() Condition {
	condition := Condition{}
	if len(o.And) > 0 {
		return *o.And[0]
	}
	return condition
}

type Condition struct {
	Operand *ConditionOperand `  @@`
	Not     *Condition        `| "NOT" @@`
	Exists  *Select           `| "EXISTS" "(" @@ ")"`
}

func (o *Condition) GetConditionOperand() ConditionOperand {
	operand := ConditionOperand{}
	if o.Operand != nil {
		return *o.Operand
	}
	return operand
}

type ConditionOperand struct {
	Operand      *Operand      `@@`
	ConditionRHS *ConditionRHS `@@?`
}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Is      *Is      `| "IS" @@`
	Between *Between `| "BETWEEN" @@`
	In      *In      `| "IN" "(" @@ ")"`
	Like    *Like    `| "LIKE" @@`
}

type Compare struct {
	Operator string         `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand       `(  @@`
	Select   *CompareSelect ` | @@ )`
}

type CompareSelect struct {
	All    bool    `(  @"ALL"`
	Any    bool    ` | @"ANY"`
	Some   bool    ` | @"SOME" )`
	Select *Select `"(" @@ ")"`
}

type Like struct {
	Not     bool     `[ @"NOT" ]`
	Operand *Operand `@@`
}

type Is struct {
	Not          bool     `[ @"NOT" ]`
	Null         bool     `( @"NULL"`
	DistinctFrom *Operand `  | "DISTINCT" "FROM" @@ )`
}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

type In struct {
	Select      *Select       `  @@`
	Expressions []*Expression `| @@ ( "," @@ )*`
}

type Operand struct {
	Summand []*Summand `@@ ( "|" "|" @@ )*`
}

type Summand struct {
	LHS *Factor `@@`
	Op  string  `[ @("+" | "-")`
	RHS *Factor `  @@ ]`
}

type Factor struct {
	LHS *Term  `@@`
	Op  string `( @("*" | "/" | "%")`
	RHS *Term  `  @@ )?`
}

type Term struct {
	Select        *Select     `  @@`
	Value         *Value      `| @@`
	SymbolRef     *SymbolRef  `| @@`
	SubExpression *Expression `| "(" @@ ")"`
}

type SymbolRef struct {
	Symbol     string        `@Ident @( "." Ident )*`
	Parameters []*Expression `( "(" @@ ( "," @@ )* ")" )?`
}

type Value struct {
	Wildcard bool     `(  @"*"`
	Number   *float64 ` | @Number`
	String   *string  ` | @String`
	Boolean  *Boolean ` | @("TRUE" | "FALSE")`
	Null     bool     ` | @"NULL"`
	Array    *Array   ` | @@ )`
}

type Array struct {
	Expressions []*Expression `"(" @@ ( "," @@ )* ")"`
}

var (
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(SELECT|FROM|TOP|DISTINCT|ALL|WHERE|GROUP|BY|HAVING|UNION|MINUS|EXCEPT|INTERSECT|ORDER|LIMIT|OFFSET|TRUE|FALSE|NULL|IS|NOT|ANY|SOME|BETWEEN|AND|OR|LIKE|AS|IN)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_\:\[\]]*`},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{`String`, `'[^']*'|"[^"]*"`},
		{`Operators`, `<>|!=|<=|>=|[-+*/%,.()=<>]`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[Expression](
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
)

func ReduceWheres(ctx context.Context, ts map[int]WhereTask, t string) (string, error) {
	var value string
	var err error
	if len(ts) == 0 {
		return "", errors.New("error Reduces wheres map is empty")
	}
	if strings.Compare(ts[0].Name, t) == 0 {
		return ts[0].Where, nil
	}
	index := 1
	whereLength := len(ts)
	lastWhere := ts[0].Where
	for index < whereLength {
		lastWhere, err = ReduceTwoWheres(lastWhere, ts[index].Where)
		if err != nil {
			return value, err
		}
		if strings.Compare(ts[index].Name, t) == 0 {
			break
		}
		index++
	}
	value = lastWhere
	return value, nil
}

func ReduceTwoWheres(s0 string, s1 string) (string, error) {
	var value strings.Builder
	sql0, err := parser.ParseString("", s0)
	if err != nil {
		return value.String(), err
	}
	sql1, err := parser.ParseString("", s1)
	if err != nil {
		return value.String(), err
	}
	statementSet := NewStringSet()
	statementMarked := NewStringSet()

	for _, orStatement := range sql0.Or {
		for _, andStatement := range orStatement.And {
			// Check Second Statement
			for _, orStatement1 := range sql1.Or {
				for _, andStatement1 := range orStatement1.And {
					similar := checkSimilary(andStatement.Operand, andStatement1.Operand)
					statement0 := WriteStruct(andStatement.Operand)
					statement1 := WriteStruct(andStatement1.Operand)
					if similar {
						statementMarked.AddString(statement0)
						statementSet.AddString(statement1)
					} else {
						if !statementMarked.ContainsString(statement0) {
							statementSet.AddString(statement0)
						}
						statementSet.AddString(statement1)
					}
				}
			}
		}
	}
	for _, statement := range statementSet.Elements() {
		value.WriteString(statement)
		value.WriteString(AndOperator)
	}
	str := value.String()
	str = strings.Trim(str, AndOperator)
	return str, nil
}

func checkSimilary(e0 *ConditionOperand, e1 *ConditionOperand) bool {
	hasSimilar := true
	var op0 *Operand
	var op1 *Operand
	var cp0 *ConditionRHS
	var cp1 *ConditionRHS
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				hasSimilar = false
			}
		}()
		op0 = e0.Operand
		op1 = e1.Operand
		cp0 = e0.ConditionRHS
		cp1 = e1.ConditionRHS
	}
	f()
	if hasSimilar == false {
		return false
	}
	// case 1
	value := reflect.DeepEqual(op0, op1) && DeepEqualCompare(cp0, cp1)
	// case 2  operand in second place
	if value == false {
		fMix := func() {
			defer func() {
				if r := recover(); r != nil {
					hasSimilar = false
				}
			}()
			op0 = e0.ConditionRHS.Compare.Operand
			op1 = e1.Operand
			cp0 = e0.ConditionRHS
			cp1 = e1.ConditionRHS
		}
		fMix()
		value = reflect.DeepEqual(op0, op1) && DeepEqualCompare(cp0, cp1)
	}
	// case 3
	if value == false {
		fMix := func() {
			defer func() {
				if r := recover(); r != nil {
					hasSimilar = false
				}
			}()
			op0 = e0.Operand
			op1 = e1.ConditionRHS.Compare.Operand
			cp0 = e0.ConditionRHS
			cp1 = e1.ConditionRHS
		}
		fMix()
		value = reflect.DeepEqual(op0, op1) && DeepEqualCompare(cp0, cp1)
	}

	return value
}

func DeepEqualCompare(c0 *ConditionRHS, c1 *ConditionRHS) bool {

	if (c0.Compare != nil && c1.Compare != nil) &&
		(c0.Compare.Operator == c1.Compare.Operator) {
		return true
	}

	if c0.In != nil && c1.In != nil {
		return true
	}

	if c0.Is != nil && c1.Is != nil {
		return true
	}

	if c0.Like != nil && c1.Like != nil {
		return true
	}

	if c0.Between != nil && c1.Between != nil {
		return true
	}

	return false
}

func ReducesFilter[T db.Scanner](wf common.Workflow[T], name string) (string, error) {
	wheres := ClausesWheres(wf)
	value := ""
	ctx := context.Background()
	value, err := ReduceWheres(ctx, wheres, name)
	if err != nil {
		return value, err
	}
	return value, nil
}

func ClausesWheres[T db.Scanner](wf common.Workflow[T]) map[int]WhereTask {
	ctx := context.Background()
	values := make(map[int]WhereTask)
	statuses, err := wf.Statuses(ctx, false)
	if err != nil {
		log.Println("Error WhereClauses Statuses")
		return values
	}
	index := 0
	str := ""
	for _, name := range statuses {
		current, _ := wf.Tasks[name]
		if index == 0 {
			values[index] = WhereTask{name, wf.Initial.String()}
		} else {
			values[index] = WhereTask{name, str}
		}
		str = current.Transition.Links[0].Query
		if strings.Compare(str, EndTask) != 0 {
			index++
		}
	}
	return values
}
