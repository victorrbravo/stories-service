package pg

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Stack is a simple stack data structure for holding reflect.Values
type Stack struct {
	items []reflect.Value
}

// Push adds a value to the stack
func (s *Stack) Push(item reflect.Value) {
	s.items = append(s.items, item)
}

func (s *Stack) Clear() {
	s.items = s.items[:0]
}

// Pop removes and returns the top value from the stack
func (s *Stack) Pop() (reflect.Value, bool) {
	if len(s.items) == 0 {
		return reflect.Value{}, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack) RemoveFront() (reflect.Value, bool) {
	if len(s.items) == 0 {
		return reflect.Value{}, false
	}
	item := s.items[0]
	s.items = s.items[1:]
	return item, true
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// WriteStruct prints the names and values of the fields of a struct using a stack
func WriteStruct(s interface{}) string {
	var buffer strings.Builder
	var v reflect.Value
	p := reflect.ValueOf(s)

	if p.Kind() == reflect.Ptr {
		if p.IsNil() {
			return buffer.String()
		}
		v = p.Elem()
	} else {
		v = reflect.ValueOf(s)
	}

	// Make sure s is a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Input is not a struct")
		return buffer.String()
	}

	stack := Stack{}
	stack.Push(v)

	operators := Stack{}
	literals := Stack{}

	for !stack.IsEmpty() {
		top, ok := stack.RemoveFront()
		if !ok {
			break
		}

		for i := 0; i < top.NumField(); i++ {
			field := top.Type().Field(i)
			fieldValue := top.Field(i)

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					continue
				}
				fieldValue = fieldValue.Elem()
			}
			switch fieldValue.Kind() {
			case reflect.Struct:
				switch field.Name {
				case "Like", "In", "Is", "Between":
					operators.Push(reflect.ValueOf(strings.ToUpper(field.Name)))
				}
				stack.Push(fieldValue)
			case reflect.Slice, reflect.Array:
				for i := 0; i < fieldValue.Len(); i++ {
					item := fieldValue.Index(i)
					if item.Kind() == reflect.Ptr {
						if !item.IsNil() {
							stack.Push(item.Elem())
						}
					}
				}
			case reflect.Int:
				buffer.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
			case reflect.String, reflect.Float64, reflect.Bool:
				switch field.Name {
				case "String":
					buffer.WriteString(fmt.Sprintf("'%v'", fieldValue.Interface()))
					literals.Push(fieldValue)
				case "Operator":
					operators.Push(fieldValue)
				case "Number", "Boolean":
					buffer.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
					literals.Push(fieldValue)
				case "Symbol":
					operator, ok := operators.Pop()
					if literals.IsEmpty() {
						buffer.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
						if ok {
							writeOperator(operator, &buffer)
						}
					} else {
						if ok {
							writeOperator(operator, &buffer)
						}
						buffer.WriteString(fmt.Sprintf("%v", fieldValue.Interface()))
						literals.Clear()
					}
				}
			default:
			}

		}
	}
	value := strings.Trim(buffer.String(), " ")
	return value
}

func writeOperator(operator reflect.Value, buffer *strings.Builder) {
	str := operator.Interface().(string)
	if isAlpha(str) {
		str = fmt.Sprintf(" %s ", str)
	}
	buffer.WriteString(str)
}

func isAlpha(s string) bool {
	alpha := regexp.MustCompile("^[a-zA-Z]*$")
	return alpha.MatchString(s)
}
