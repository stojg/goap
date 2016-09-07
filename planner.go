// Package goap is a Goal Orientated Action Planner
// it was inspired by the source for of sploreg at
// https://github.com/sploreg/goap/blob/d1cea0728fb4733266affea8049da1e373d618f7/Assets/Standard%20Assets/Scripts/AI/GOAP/GoapPlanner.cs
package goap

type Agent interface{}

type Action interface {
	Reset()
	CheckProceduralPrecondition(agent Agent) bool
	Preconditions() KeyValuePairs
	Effects() KeyValuePairs
	Cost() float64
}

type Comparable interface {
	Equals(Comparable) bool
	Key() string
	Value() interface{}
}

type KeyValuePairs map[string]Comparable

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

	//List<Node> leaves = new List<Node>();
	var leaves []*node
	//
	//// build graph
	start := newNode(nil, 0, worldState, nil)
	success := buildGraph(start, leaves, usableActions, goal)

	if !success {
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

	// go through the end node and work up to it's parents
	for n != nil {
		if n.action != nil {
			// insert action in front
			result = append([]Action{n.action}, result...)
		}
		n = n.parent
	}

	if len(result) > 0 {
		return result
	}
	return nil
}

// buildGraph returns true if at least one solution was found. The possible paths are stored in the
// leaves list. Each leaf has a 'runningCost' value where the lowest cost will be the best action
// sequence.
func buildGraph(parent *node, leaves []*node, usableActions []Action, goal KeyValuePairs) bool {
	foundOne := false

	// go through each action available at this node and see if we can use it here
	for _, action := range usableActions {
		// if the parent state has the conditions for this action's preconditions, we can use it here
		if inState(action.Preconditions(), parent.state) {

			// apply the action's effects to the parent state
			var currentState = populateState(parent.state, action.Effects())
			node := newNode(parent, parent.runningCost+action.Cost(), currentState, action)
			if inState(goal, currentState) {
				// we found a solution!
				// @todo should the node be added first in the list?
				leaves = append(leaves, node)
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
	}
	return foundOne
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

// apply the statechange to the currect state
func populateState(currentState KeyValuePairs, stateChange KeyValuePairs) KeyValuePairs {
	var state KeyValuePairs

	// copy the KVPs over as new objects
	for _, s := range currentState {
		state[s.Key()] = s
	}

	// if the key exists in the current state, update or add the Value
	for _, change := range stateChange {
		state[change.Key()] = change
	}
	return state
}

// Check that all items in 'test' are in 'state'. If just one does not match or is not there then
// this returns false.
func inState(test KeyValuePairs, state KeyValuePairs) bool {
	for _, t := range test {

		match := false

		for _, s := range state {
			// find at least one match in the world state
			if s.Equals(t) {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}
	return true
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
