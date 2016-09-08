package goap

import "testing"

type testAgent struct{}

type eatFoodAction struct {
	name string
}

func (a *eatFoodAction) Reset() {

}
func (a *eatFoodAction) CheckProceduralPrecondition(agent Agent) bool {
	return true
}
func (a *eatFoodAction) Preconditions() KeyValuePairs {
	return make(KeyValuePairs, 0)
}
func (a *eatFoodAction) Effects() KeyValuePairs {
	return make(KeyValuePairs, 0)
}
func (a *eatFoodAction) Cost() float64 {
	return 1
}

type state struct {
	key   string
	value int
}

func Test_populateState(t *testing.T) {

	worldState := make(KeyValuePairs, 0)
	worldState["food"] = 1
	worldState["temperature"] = 10

	changes := make(KeyValuePairs, 0)
	changes["food"] = 2

	result := populateState(worldState, changes)

	if len(result) != len(worldState) {
		t.Errorf("result should have the same # of entries as world state")
	}

	if _, ok := result["food"]; !ok {
		t.Logf("%s", result)
		t.Errorf("could not find 'food' state")
		return
	}

	if result["food"] != 2 {
		t.Errorf("food state was not changed, expected 2, got %d", result["food"])
	}

	if worldState["food"] != 1 {
		t.Errorf("worldState failed to be treated as an immutable")
	}

	if result["temperature"] != 10 {
		t.Errorf("unrelated state was changed, temperature is ", result["temperature"])
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
		t.Errorf("expected that inState would be true")
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
		t.Errorf("expected that inState would be false")
	}
}

//func TestPlan(t *testing.T) {
//
//	agent := &testAgent{}
//
//	worldState := make(KeyValuePairs, 0)
//	worldState["food"] = &state{
//		value: 1,
//	}
//
//	goal := make(KeyValuePairs, 0)
//	goal["food"] = &state{
//		key: "food",
//		value: 0,
//	}
//	actions := []Action{
//		&eatFoodAction{ name: "first", },
//	}
//
//	res := Plan(agent, actions, worldState, goal)
//
//	if res == nil {
//		t.Errorf("Did not find a plan")
//	}
//}
