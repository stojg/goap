package goap

import (
	"testing"
)

type testAgent struct{}

type state struct {
	key   string
	value int
}

func TestPlan1(t *testing.T) {

	agent := &testAgent{}

	worldState := make(KeyValuePairs, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false

	getFood := newTestAction("getFood", 8)
	getFood.AddPrecondition("hasFood", false)
	getFood.AddEffect("hasFood", true)

	eat := newTestAction("eat", 4)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	sleep := newTestAction("sleep", 4)
	sleep.AddPrecondition("isTired", true)
	sleep.AddEffect("isTired", false)

	availableActions := []Action{getFood, eat, sleep}

	goal := make(KeyValuePairs, 0)
	goal["isFull"] = true

	actionList := Plan(agent, availableActions, worldState, goal)

	if actionList == nil {
		t.Error("Expected to get a plan, got no plan")
	}

	expectedActions := 2
	if len(actionList) != expectedActions {
		t.Errorf("There should be %d actions in the plan, got %d", expectedActions, len(actionList))
		t.Errorf("planned actions: %+v", actionList)
	}

	if actionList[0].String() != "getFood" {
		t.Errorf("expected first action to be 'getFood', but got %s", actionList[0])
	}

	if actionList[1].String() != "eat" {
		t.Errorf("expected second action to be 'eat', but got %s", actionList[1])
	}
}

func TestPlan2(t *testing.T) {

	agent := &testAgent{}

	worldState := make(KeyValuePairs, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false

	getFood := newTestAction("getFood", 8)
	getFood.AddPrecondition("hasFood", false)
	getFood.AddEffect("hasFood", true)

	// test that the planner finds the cheapest way to the same goal
	prayForFood := newTestAction("prayForFood", 6)
	prayForFood.AddPrecondition("hasFood", false)
	prayForFood.AddEffect("hasFood", true)

	eat := newTestAction("eat", 4)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	sleep := newTestAction("sleep", 4)
	sleep.AddPrecondition("isTired", true)
	sleep.AddEffect("isTired", false)

	availableActions := []Action{getFood, prayForFood, eat, sleep}

	goal := make(KeyValuePairs, 0)
	goal["isFull"] = true

	actionList := Plan(agent, availableActions, worldState, goal)

	if actionList == nil {
		t.Error("Expected to get a plan, got no plan")
	}

	expectedActions := 2
	if len(actionList) != expectedActions {
		t.Errorf("There should be %d actions in the plan, got %d", expectedActions, len(actionList))
		t.Errorf("planned actions: %+v", actionList)
	}

	if actionList[0].String() != "prayForFood" {
		t.Errorf("expected first action to be 'prayForFood', but got %s", actionList[0])
	}

	if actionList[1].String() != "eat" {
		t.Errorf("expected second action to be 'eat', but got %s", actionList[1])
	}
}

func TestPlan3_failed(t *testing.T) {

	agent := &testAgent{}

	worldState := make(KeyValuePairs, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false

	getFood := newTestAction("getFood", 8)
	getFood.AddPrecondition("hasFood", false)
	getFood.AddEffect("hasFood", true)

	eat := newTestAction("eat", 4)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	sleep := newTestAction("sleep", 4)
	sleep.AddPrecondition("isTired", true)
	sleep.AddEffect("isTired", false)

	availableActions := []Action{getFood, eat, sleep}

	// there are no actions that can fulfill this goal
	goal := make(KeyValuePairs, 0)
	goal["isWarm"] = true

	actionList := Plan(agent, availableActions, worldState, goal)

	if actionList != nil {
		t.Error("Expected the planning to fail, but it didn't")
	}
}

func Test_buildGraph(t *testing.T) {

	eat := newTestAction("eat", 4)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	eatSlowly := newTestAction("eatSlowly", 8)
	eatSlowly.AddPrecondition("hasFood", true)
	eatSlowly.AddPrecondition("isFull", false)
	eatSlowly.AddEffect("isFull", true)

	hide := newTestAction("hide", 2)
	hide.AddPrecondition("isHurt", true)
	hide.AddEffect("isHidden", true)

	usableActions := []Action{eat, eatSlowly, hide}

	goal := make(KeyValuePairs, 0)
	goal["isFull"] = true

	worldState := make(KeyValuePairs, 0)
	worldState["hasFood"] = true
	worldState["isFull"] = false

	start := newNode(nil, 0, worldState, nil)

	var leaves []*node
	found := buildGraph(start, &leaves, usableActions, goal)

	if !found {
		t.Error("expected to find a plan")
	}

	if len(leaves) < 1 {
		t.Errorf("expected at least one item in the plan, got %d", len(leaves))
	}
}

func Test_inState_true(t *testing.T) {

	test := make(KeyValuePairs, 0)
	test["food"] = 2

	state := make(KeyValuePairs, 0)
	state["food"] = 2
	state["temperature"] = 10

	actual := inState(test, state)
	if !actual {
		t.Error("expected that inState would be true")
	}
}

func Test_inState_false(t *testing.T) {

	test := make(KeyValuePairs, 0)
	test["food"] = 1

	state := make(KeyValuePairs, 0)
	state["food"] = 2
	state["temperature"] = 10

	actual := inState(test, state)
	if actual {
		t.Error("expected that inState would be false")
	}
}

func Test_inState_dont_exists(t *testing.T) {

	test := make(KeyValuePairs, 0)
	test["isHurt"] = true

	state := make(KeyValuePairs, 0)
	state["hasFood"] = true
	state["isFull"] = false

	actual := inState(test, state)
	if actual {
		t.Error("expected that inState would be false")
	}
}

func Test_actionSubset(t *testing.T) {
	eat := &testAction{name: "eat"}
	drop := &testAction{name: "drop"}
	hide := &testAction{name: "hide"}

	actions := []Action{eat, drop, hide}

	result := actionSubset(actions, drop)

	if len(result) != 2 {
		t.Errorf("expected one action to be removed, got %d", len(result))
	}

	if result[0] != eat {
		t.Error("expected eat action")
	}

	if result[1] == drop {
		t.Error("didnt expected drop action")
	}

	if result[1] != hide {
		t.Error("did expected hide actiond")
	}
}

func Test_populateState(t *testing.T) {

	worldState := make(KeyValuePairs, 0)
	worldState["food"] = 1
	worldState["temperature"] = 10

	changes := make(KeyValuePairs, 0)
	changes["food"] = 2

	result := populateState(worldState, changes)

	if len(result) != len(worldState) {
		t.Error("result should have the same # of entries as world state")
	}

	if _, ok := result["food"]; !ok {
		t.Logf("%s", result)
		t.Error("could not find 'food' state")
		return
	}

	if result["food"] != 2 {
		t.Errorf("food state was not changed, expected 2, got %d", result["food"])
	}

	if worldState["food"] != 1 {
		t.Error("worldState failed to be treated as an immutable")
	}

	if result["temperature"] != 10 {
		t.Error("unrelated state was changed, temperature is ", result["temperature"])
	}
}

func newTestAction(name string, cost float64) *testAction {
	return &testAction{
		name:          name,
		cost:          cost,
		preconditions: make(KeyValuePairs, 0),
		effects:       make(KeyValuePairs, 0),
	}
}

type testAction struct {
	name          string
	cost          float64
	preconditions KeyValuePairs
	effects       KeyValuePairs
}

func (a *testAction) Reset() {}

func (a *testAction) CheckProceduralPrecondition(agent Agent) bool {
	return true
}
func (a *testAction) Preconditions() KeyValuePairs {
	return a.preconditions
}
func (a *testAction) Effects() KeyValuePairs {
	return a.effects
}
func (a *testAction) Cost() float64 {
	return a.cost
}
func (a *testAction) AddPrecondition(key string, value interface{}) {
	a.preconditions[key] = value
}
func (a *testAction) RemovePrecondition(key string) {
	delete(a.preconditions, key)
}
func (a *testAction) AddEffect(key string, value interface{}) {
	a.effects[key] = value
}

func (a *testAction) RemoveEffect(key string) {
	delete(a.effects, key)
}

func (a *testAction) IsDone() bool {
	return false
}

func (a *testAction) IsInRange() bool {
	return true
}

func (a *testAction) Perform(agent Agent) bool {
	return true
}

func (a *testAction) RequiresInRange() bool {
	return false
}

func (a *testAction) String() string {
	return a.name
}
