package goap

// Actionable is the interface that describes what the planner and
type Actionable interface {

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
	CheckProceduralPrecondition(Agent) bool

	// Run the action.
	// Returns True if the action performed successfully or false
	// if something happened and it can no longer perform. In this case
	// the action queue should clear out and the goal cannot be reached.
	Perform(Agent) bool

	// Does this action need to be within range of a target game object?
	// If not then the moveTo state will not need to run for this action.
	RequiresInRange() bool
	IsInRange() bool
	SetInRange()
	SetTarget(interface{})
	Target() interface{}

	AddPrecondition(key string, value interface{})
	RemovePrecondition(key string)
	Preconditions() StateList

	AddEffect(key string, value interface{})
	RemoveEffect(key string)
	Effects() StateList

	String() string
}

// NewAction create a new base Action
func NewAction(name string, cost float64) Action {
	return Action{
		name:          name,
		preconditions: make(StateList, 0),
		effects:       make(StateList, 0),
		cost:          cost,
	}
}

type Action struct {
	name          string
	preconditions StateList
	effects       StateList
	cost          float64
	isDone        bool
	inRange       bool
	target        interface{}
}

func (a *Action) Reset() {
	a.isDone = false
	a.inRange = false
	a.target = nil
}

func (a *Action) AddPrecondition(key string, value interface{}) {
	a.preconditions[key] = value
}

func (a *Action) RemovePrecondition(key string) {
	delete(a.preconditions, key)
}

func (a *Action) Preconditions() StateList {
	return a.preconditions
}

func (a *Action) AddEffect(key string, value interface{}) {
	a.effects[key] = value
}

func (a *Action) RemoveEffect(key string) {
	delete(a.effects, key)
}

func (a *Action) Effects() StateList {
	return a.effects
}

func (a *Action) Cost() float64 {
	return a.cost
}

func (a *Action) IsInRange() bool {
	return a.inRange
}

func (a *Action) SetInRange() {
	a.inRange = true
}

func (a *Action) IsDone() bool {
	return a.isDone
}

func (a *Action) SetTarget(t interface{}) {
	a.target = t
}

func (a *Action) Target() interface{} {
	return a.target
}

func (a *Action) String() string {
	return a.name
}
