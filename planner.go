// Package goap is a Goal Orientated Action Planner
// it was inspired by the source for of sploreg at
// https://github.com/sploreg/goap/blob/d1cea0728fb4733266affea8049da1e373d618f7/Assets/Standard%20Assets/Scripts/AI/GOAP/GoapPlanner.cs
package goap

type Agent interface{}

type Action interface {

	// he cost of performing the action.
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

type Comparable interface {
	Value() interface{}
}

type KeyValuePairs map[string]interface{}

// Plan what sequence of actions can fulfill the goal. Returns null if a plan could not be found, or
// a list of the actions that must be performed, in order, to fulfill the goal.
func Plan(agent Agent, availableActions []Action, worldState KeyValuePairs, goal KeyValuePairs) []Action {

	// reset the actions so we can start fresh with them
	for i := range availableActions {
		availableActions[i].Reset()
	}

	// check what actions can run
	var usableActions []Action
	for _, a := range availableActions {
		if a.CheckProceduralPrecondition(agent) {
			usableActions = append(usableActions, a)
		}
	}

	// we now have all actions that can run, stored in usableActions
	// build up the tree and record the leaf nodes that provide a solution to the goal.
	var leaves []*node
	start := newNode(nil, 0, worldState, nil)
	if !buildGraph(start, &leaves, usableActions, goal) {
		return nil
	}

	// get the cheapest leaf
	var cheapest *node
	for _, leaf := range leaves {
		if cheapest == nil {
			cheapest = leaf
		} else if leaf.runningCost < cheapest.runningCost {
			cheapest = leaf
		}
	}

	var result []Action
	n := cheapest

	// go through the end node and work up to through it's parents
	for n != nil {
		if n.action != nil {
			// insert action in front
			result = append([]Action{n.action}, result...)
		}
		n = n.parent
	}
	return result
}

// buildGraph returns true if at least one solution was found. The possible paths are stored in the
// leaves list. Each leaf has a 'runningCost' value where the lowest cost will be the best action
// sequence.
func buildGraph(parent *node, leaves *[]*node, usableActions []Action, goal KeyValuePairs) bool {
	foundOne := false

	// go through each action available at this node and see if we can use it here
	for _, action := range usableActions {

		// if the parent state has the conditions for this action's preconditions, we can use it here
		if !inState(action.Preconditions(), parent.state) {
			continue
		}

		// apply the action's effects to the parent state
		var currentState = populateState(parent.state, action.Effects())
		node := newNode(parent, parent.runningCost+action.Cost(), currentState, action)

		if inState(goal, currentState) {
			// we found a solution!
			*leaves = append(*leaves, node)
			foundOne = true
		} else {
			// not at a solution yet, so test all the remaining actions and branch out the tree
			subset := actionSubset(usableActions, action)
			found := buildGraph(node, leaves, subset, goal)
			if found {
				foundOne = true
			}
		}
	}
	return foundOne
}

// Check that all items in 'test' are in 'state'. If just one does not match or is not there then
// this returns false.
func inState(test KeyValuePairs, state KeyValuePairs) bool {
	for testKey, testVal := range test {
		match := false
		if stateVal, found := state[testKey]; found {
			if stateVal == testVal {
				match = true
			}
		}
		if !match {
			return false
		}
	}
	return true
}

// apply the state change to the current state
func populateState(currentState KeyValuePairs, stateChange KeyValuePairs) KeyValuePairs {
	state := make(KeyValuePairs, 0)

	// copy the KVPs over as new objects
	for key, s := range currentState {
		state[key] = s
	}

	// if the key exists in the current state, update or add the Value
	for key, change := range stateChange {
		state[key] = change
	}
	return state
}

func actionSubset(actions []Action, removeMe Action) []Action {
	var subset []Action
	for _, a := range actions {
		if a != removeMe {
			subset = append(subset, a)
		}
	}
	return subset
}

// Node is used for building up the graph and holding the running costs of actions.
type node struct {
	parent      *node
	runningCost float64
	state       KeyValuePairs
	action      Action
}

func newNode(parent *node, runningCost float64, state KeyValuePairs, action Action) *node {
	return &node{
		parent:      parent,
		runningCost: runningCost,
		state:       state,
		action:      action,
	}
}