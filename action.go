package goap

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

	AddPrecondition(key string, value interface{})
	RemovePrecondition(key string)
	Preconditions() KeyValuePairs

	AddEffect(key string, value interface{})
	RemoveEffect(key string)
	Effects() KeyValuePairs

	String() string
}
