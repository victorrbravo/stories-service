package utils

import (
	"errors"
)

// Stack represents a stack data structure
type Stack struct {
	elements []string
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{elements: []string{}}
}

// Push adds an element to the top of the stack
func (s *Stack) Push(element string) {
	s.elements = append(s.elements, element)
}

// Pop removes and returns the top element of the stack
// It returns an error if the stack is empty
func (s *Stack) Pop() (string, error) {
	if len(s.elements) == 0 {
		return "", errors.New("stack is empty")
	}
	topIndex := len(s.elements) - 1
	topElement := s.elements[topIndex]
	s.elements = s.elements[:topIndex]
	return topElement, nil
}

// Peek returns the top element of the stack without removing it
// It returns an error if the stack is empty
func (s *Stack) Peek() (string, error) {
	if len(s.elements) == 0 {
		return "", errors.New("stack is empty")
	}
	return s.elements[len(s.elements)-1], nil
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return len(s.elements)
}
