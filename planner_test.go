package goap

import (
	"testing"
)

func TestPlan1(t *testing.T) {

	agent := &DefaultAgent{}

	actions := []Action{findFood(), eatAction(), sleepAction()}

	currentState := make(StateList)
	currentState.Is(Hungry).Dont(HaveFood)

	goal := make(StateList)
	goal.Isnt(Hungry)

	actionList := Plan(agent, actions, currentState, goal)

	if actionList == nil {
		t.Error("Expected to get a plan, got no plan")
	}

	expectedActions := 2
	if len(actionList) != expectedActions {
		t.Errorf("There should be %d actions in the plan, got %d", expectedActions, len(actionList))
		t.Errorf("planned actions: %+v", actionList)
		return
	}

	if actionList[0].String() != "getFood" {
		t.Errorf("expected first action to be 'getFood', but got %s", actionList[0])
		return
	}

	if actionList[1].String() != "eat" {
		t.Errorf("expected second action to be 'eat', but got %s", actionList[1])
		return
	}
}

func TestPlan2(t *testing.T) {

	agent := &DefaultAgent{}

	// test that the planner finds the cheapest way to the same goal
	prayForFood := newTestAction("prayForFood", 6, false)
	prayForFood.AddEffect(HaveFood)
	prayForFood.AddPrecondition(Dont(HaveFood))

	actions := []Action{findFood(), prayForFood, eatAction(), sleepAction()}

	currentState := make(StateList)
	currentState.Is(Hungry).Dont(HaveFood)

	goal := make(StateList)
	goal.Isnt(Hungry)

	actionList := Plan(agent, actions, currentState, goal)

	if actionList == nil {
		t.Error("Expected to get a plan, got no plan")
	}

	expectedActions := 2
	if len(actionList) != expectedActions {
		t.Errorf("There should be %d actions in the plan, got %d", expectedActions, len(actionList))
		t.Errorf("planned actions: %+v", actionList)
		return
	}

	if actionList[0].String() != "prayForFood" {
		t.Errorf("expected first action to be 'prayForFood', but got %s", actionList[0])
		return
	}

	if actionList[1].String() != "eat" {
		t.Errorf("expected second action to be 'eat', but got %s", actionList[1])
		return
	}
}

func TestPlan_failed(t *testing.T) {

	agent := &DefaultAgent{}

	actions := []Action{findFood(), eatAction(), sleepAction()}

	// there are no actions that can fulfill this goal
	goal := make(StateList)
	goal.Add(State{"isWarm", true})

	currentState := make(StateList)
	currentState.Is(Hungry).Dont(HaveFood)

	actionList := Plan(agent, actions, currentState, goal)

	if actionList != nil {
		t.Error("Expected the planning to fail, but it didn't")
	}
}

func Test_buildGraph(t *testing.T) {
	eatSlowly := newTestAction("eatSlowly", 8, false)
	eatSlowly.AddEffect(Isnt(Hungry), Dont(HaveFood))
	eatSlowly.AddPrecondition(Hungry, HaveFood)

	hide := newTestAction("hide", 2, false)
	hide.AddPrecondition(IsHurt)
	hide.AddEffect(IsHidden)

	actions := []Action{eatAction(), eatSlowly, hide}

	currentState := make(StateList)
	currentState.Add(HaveFood).Is(Hungry)

	goal := make(StateList)
	goal.Isnt(Hungry)

	start := newNode(nil, 0, currentState, nil)

	var leaves []*node
	found := buildGraph(start, &leaves, actions, goal)

	if !found {
		t.Error("expected to find a plan")
		return
	}

	if len(leaves) < 1 {
		t.Errorf("expected at least one item in the plan, got %d", len(leaves))
	}
}

func Test_inState_true(t *testing.T) {

	test := make(StateList)
	test["food"] = true

	state := make(StateList)
	state["food"] = true
	state["temperature"] = false

	actual := inState(test, state)
	if !actual {
		t.Error("expected that inState would be true")
	}
}

func Test_inState_false(t *testing.T) {

	test := make(StateList)
	test["food"] = true

	state := make(StateList)
	state["food"] = false
	state["temperature"] = true

	actual := inState(test, state)
	if actual {
		t.Error("expected that inState would be false")
	}
}

func Test_inState_dont_exists(t *testing.T) {

	test := make(StateList)
	test["isHurt"] = true

	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	actual := inState(test, state)
	if actual {
		t.Error("expected that inState would be false")
	}
}

func Test_actionSubset(t *testing.T) {
	eat := newTestAction("eat", 2, false)
	drop := newTestAction("drop", 2, false)
	hide := newTestAction("hide", 2, false)

	actions := []Action{eat, drop, hide}

	result := actionSubset(actions, drop)

	if len(result) != 2 {
		t.Errorf("expected one action to be removed, got %d", len(result))
	}

	if result[0] != eat {
		t.Error("expected eat action")
		return
	}

	if result[1] == drop {
		t.Error("didnt expected drop action")
		return
	}

	if result[1] != hide {
		t.Error("did expected hide actiond")
		return
	}
}

func Test_populateState(t *testing.T) {

	currentState := make(StateList)
	currentState["food"] = false
	currentState["temperature"] = true

	changes := make(StateList)
	changes["food"] = true

	result := populateState(currentState, changes)

	if len(result) != len(currentState) {
		t.Error("result should have the same # of entries as world state")
	}

	if _, ok := result["food"]; !ok {
		t.Logf("%s", result.String())
		t.Error("could not find 'food' state")
		return
	}

	if !result["food"] {
		t.Errorf("food state was not changed, expected true, got %t", result["food"])
	}

	if currentState["food"] {
		t.Error("currentState failed to be treated as an immutable")
	}

	if !result["temperature"] {
		t.Error("unrelated state was changed, temperature is ", result["temperature"])
	}
}

func newTestAction(name string, cost float64, requiresInRange bool) *testAction {
	action := &testAction{
		DefaultAction: NewAction(name, cost),
	}
	if !requiresInRange {
		action.inRange = true
	}
	return action
}

type testAction struct {
	DefaultAction
	inRange bool
}

func (a *testAction) Perform(agent Agent) bool {
	return true
}

func (a *testAction) InRange(agent Agent) bool {
	return true
}
