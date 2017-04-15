package goap

// Action is the interface that describes what the planner and
type Action interface {
	// The cost of performing the action.
	// Figure out a weight that suits the action.
	// Changing it will affect what actions are chosen during planning.
	Cost() float64

	// Reset any variables that need to be reset before planning happens again.
	Reset()

	// Is the action done
	IsDone() bool

	// Procedurally check if this action can run. Not all actions will need
	// this, but some might.
	CheckContextPrecondition(Agent) bool

	// Run the action.
	// Returns True if the action performed successfully or false
	// if something happened and it can no longer perform. In this case
	// the action queue should clear out and the goal cannot be reached.
	Perform(Agent) bool

	// Does this action need to be within range of a target game object?
	// If not then the moveTo state will not need to run for this action.
	InRange(Agent) bool

	SetTarget(interface{})
	Target() interface{}

	AddPrecondition(...State)
	RemovePrecondition(key string)
	Preconditions() StateList

	AddEffect(...State)
	RemoveEffect(key string)
	Effects() StateList

	String() string
}

// NewAction create a new base DefaultAction
func NewAction(name string, cost float64) DefaultAction {
	return DefaultAction{
		name:          name,
		preconditions: make(StateList),
		effects:       make(StateList),
		cost:          cost,
	}
}

type DefaultAction struct {
	name            string
	preconditions   StateList
	effects         StateList
	cost            float64
	Done            bool
	requiresInRange bool
	target          interface{}
}

func (a *DefaultAction) Reset() {
	a.Done = false
	a.target = nil
}

func (a *DefaultAction) AddPrecondition(states ...State) {
	for _, state := range states {
		a.preconditions[state.Name] = state.Value
	}
}

func (a *DefaultAction) RemovePrecondition(key string) {
	delete(a.preconditions, key)
}

func (a *DefaultAction) Preconditions() StateList {
	return a.preconditions
}

func (a *DefaultAction) AddEffect(states ...State) {
	for _, state := range states {
		a.effects[state.Name] = state.Value
	}
}

func (a *DefaultAction) RemoveEffect(key string) {
	delete(a.effects, key)
}

func (a *DefaultAction) Effects() StateList {
	return a.effects
}

func (a *DefaultAction) Cost() float64 {
	return a.cost
}

func (a *DefaultAction) IsDone() bool {
	return a.Done
}

func (a *DefaultAction) SetTarget(t interface{}) {
	a.target = t
}

func (a *DefaultAction) Target() interface{} {
	return a.target
}

func (a *DefaultAction) CheckContextPrecondition(agent Agent) bool {
	return true
}

func (a *DefaultAction) String() string {
	return a.name
}
