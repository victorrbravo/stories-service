package pg

// StringSet is a set implementation specifically for strings
type StringSet struct {
	data map[string]bool
}

// NewStringSet creates a new empty set for strings
func NewStringSet() *StringSet {
	return &StringSet{data: make(map[string]bool)}
}

// AddString adds a string element to the set
func (s *StringSet) AddString(element string) {
	s.data[element] = true
}

// RemoveString removes a string element from the set
func (s *StringSet) RemoveString(element string) {
	delete(s.data, element)
}

// ContainsString checks if the set contains a specific string element
func (s *StringSet) ContainsString(element string) bool {
	return s.data[element]
}

// Size returns the number of string elements in the set
func (s *StringSet) Size() int {
	return len(s.data)
}

// Elements returns a slice of all string elements in the set
func (s *StringSet) Elements() []string {
	elements := make([]string, 0, len(s.data))
	for element := range s.data {
		elements = append(elements, element)
	}
	return elements
}
