package handlers_utils

import (
	"net/http"
	"regexp"
)

type RouteNode struct {
	child   *RouteNode
	next    *RouteNode
	key     string
	pattern string
	handler http.Handler
}

func newRouteNode(name string) *RouteNode {
	var regex *regexp.Regexp = regexp.MustCompile("{([^}]+)}")
	var matches [][]string = regex.FindAllStringSubmatch(name, -1)
	if len(matches) != 0 {
		return &RouteNode{
			key:     name,
			pattern: matches[0][1],
		}
	} else {
		return &RouteNode{
			key: name,
		}
	}
}

func (node *RouteNode) GetMethods() []string {
	var array []string = make([]string, 0)
	var child *RouteNode
	for child = node.child; child != nil; child = child.next {
		array = append(array, child.GetMethods()...)
		switch child.key {
		case http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut:
			array = append(array, child.key)
			break
		}
	}
	return array
}

func (node *RouteNode) FindParametrizedSibling(name string) *RouteNode {
	for sibling := node; sibling != nil; sibling = sibling.next {
		if sibling.pattern != "" {
			return sibling
		}
	}
	return nil
}

func (node *RouteNode) FindMatchingSibling(name string) *RouteNode {
	for sibling := node; sibling != nil; sibling = sibling.next {
		if sibling.key == name {
			return sibling
		}
	}
	return nil
}

func (node *RouteNode) FindMatchingChild(name string) *RouteNode {
	var child *RouteNode
	var found *RouteNode
	child = node.child
	if child == nil {
		return nil
	}
	found = child.FindMatchingSibling(name)
	if found == nil {
		return child.FindParametrizedSibling(name)
	} else {
		return found
	}
}

func (node *RouteNode) AddFragment(name string) (*RouteNode, error) {
	var child *RouteNode
	child = node.child
	// If there's no child yet, then add it directly and that's it
	if child == nil {
		node.child = newRouteNode(name)
		// Return the child so that the caller can keep
		// adding children (or siblings)
		return node.child, nil
	} else {
		// Else, there's a child find one that matches our fragment
		var found *RouteNode
		found = child.FindMatchingSibling(name)
		// If none is found, we need to add one
		if found == nil {
			var last *RouteNode
			// Find the last child
			for last = child; last.next != nil; last = last.next {
			}
			// Append to it
			last.next = newRouteNode(name)
			return last.next, nil
		} else {

			// Otherwise, just return the existing node
			return found, nil
		}
	}
}

func (node *RouteNode) String() string {
	if node == nil {
		return "nil"
	}
	return node.key
}
