package goap

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
