package goap

import (
	"fmt"
	"strings"
)

type StateList map[string]bool

func (s *StateList) Add(n State) *StateList {
	(*s)[n.Name] = n.Value
	return s
}

func (s *StateList) Is(n State) *StateList {
	(*s)[n.Name] = n.Value
	return s
}

func (s *StateList) Isnt(n State) *StateList {
	(*s)[n.Name] = !n.Value
	return s
}

func (s *StateList) Dont(n State) *StateList {
	(*s)[n.Name] = !n.Value
	return s
}

func (s StateList) Query(n State) bool {
	if v, ok := s[n.Name]; ok {
		return v == n.Value
	}
	return false
}

func (s *StateList) String() string {
	var res []string
	for k, v := range *s {
		res = append(res, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(res, ", ")
}

type State struct {
	Name  string
	Value bool
}

func Isnt(s State) State {
	return State{
		Name:  s.Name,
		Value: !s.Value,
	}
}

func Dont(s State) State {
	return State{
		Name:  s.Name,
		Value: !s.Value,
	}
}
