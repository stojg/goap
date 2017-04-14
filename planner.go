// Package goap is a Goal Orientated Action Planner just mainly for game programming purposes
// Inspired by https://github.com/sploreg/goap/
package goap

// Plan what sequence of actions can fulfill the goal. Returns null if a plan could not be found, or
// a list of the actions that must be performed, in order, to fulfill the goal.
//
func Plan(agent Agent, availableActions []Actionable, worldState StateList, goal StateList) []Actionable {

	// check what actions can run
	var usableActions []Actionable
	for _, action := range availableActions {
		// reset the actions so we can start fresh with them
		action.Reset()
		if action.SetAgent(agent) {
			usableActions = append(usableActions, action)
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

	// go through the end node and work up to through it's parents
	var result []Actionable
	for n := cheapest; n != nil; n = n.parent {
		if n.action != nil {
			// insert action in front
			result = append([]Actionable{n.action}, result...)
		}
	}
	return result
}

// buildGraph returns true if at least one solution was found. The possible paths are stored in the
// leaves list. Each leaf has a 'runningCost' value where the lowest cost will be the best action
// sequence.
func buildGraph(parent *node, leaves *[]*node, usableActions []Actionable, goal StateList) bool {
	foundOne := false

	// go through each action available at this node and see if we can use it here
	for _, action := range usableActions {

		// if the parent state has the conditions for this action's preconditions, we can use it here
		if !inState(action.Preconditions(), parent.state) {
			continue
		}

		// apply the action's effects to the parent state
		currentState := populateState(parent.state, action.Effects())
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
func inState(test StateList, state StateList) bool {
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
func populateState(currentState StateList, stateChange StateList) StateList {
	state := make(StateList)

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

func actionSubset(actions []Actionable, removeMe Actionable) []Actionable {
	var subset []Actionable
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
	state       StateList
	action      Actionable
}

func newNode(parent *node, runningCost float64, state StateList, action Actionable) *node {
	return &node{
		parent:      parent,
		runningCost: runningCost,
		state:       state,
		action:      action,
	}
}
