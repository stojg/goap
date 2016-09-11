package goap_test

import (
	"fmt"
	. "github.com/stojg/goap"
)

func NewExampleAgent(dataProvider DataProvider, actions []Actionable) *ExampleAgent {
	agent := &ExampleAgent{
		GoapAgent: NewGoapAgent(dataProvider, actions),
	}
	return agent
}

type ExampleAgent struct {
	*GoapAgent
	frame int
}

func (a *ExampleAgent) Update() {
	a.frame++
	if a.Debug {
		fmt.Printf("#%d\n", a.frame)
	}
	a.StateMachine().Update(a)
}

func ExamplePlan() {
	getFood := newGetFoodAction(8)
	getFood.AddPrecondition("hasFood", false)
	getFood.AddEffect("hasFood", true)

	eat := newEatAction(4)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	sleep := newSleepAction(4)
	sleep.AddPrecondition("isTired", true)
	sleep.AddEffect("isTired", false)

	actions := []Actionable{getFood, eat, sleep}
	provider := &dataProvider{}
	agent := NewExampleAgent(provider, actions)
	agent.Debug = true

	// 1. idle state, will do planning
	agent.Update()

	// 2. perform action getFood, but discovers that will need to move
	agent.Update()

	// 3. Move to food, it' instantly succeeds
	agent.Update()

	// 4. We have moved and food is in range
	provider.moveResult = true
	//getFood.inRange = true
	agent.Update()

	// 5. mark the getFoodAction as done
	agent.Update()

	// 6. time to eat that food
	agent.Update()
	//eat.isDone = true

	// 7. We should be done here
	agent.Update()

	// Output:
	// #1
	// idle - is planning
	// Planning success to goal map[isFull:true] with actions [getFood eat]
	// #2
	// doAction - scheduling moveTo getFood
	// #3
	// moveTo - MoveAgent(getFood)
	// #4
	// moveTo - MoveAgent(getFood)
	// moveTo - done
	// #5
	// doAction - getFood.Perform()
	// #6
	// doAction - action getFood is done
	// doAction - eat.Perform()
	// #7
	// doAction - action eat is done
	// all actions finished
}

func newGetFoodAction(cost float64) *getFoodAction {
	return &getFoodAction{
		Action: NewAction("getFood", cost),
	}
}

type getFoodAction struct {
	Action
	inRange bool
	hasFood bool
}

func (a *getFoodAction) Reset() {}

func (a *getFoodAction) CheckProceduralPrecondition(agent Agent) bool {
	a.SetTarget([]int{10, 0, 200})
	return true
}

func (a *getFoodAction) RequiresInRange() bool {
	return true
}

func (a *getFoodAction) IsInRange() bool {
	// first time we call this we are not in range, but next time yes
	if !a.inRange {
		a.inRange = true
		return false
	}
	return true
}

// Perform will
func (a *getFoodAction) Perform(agent Agent) bool {
	a.hasFood = true
	return true
}

func (a *getFoodAction) IsDone() bool {
	return a.hasFood
}

func newEatAction(cost float64) *eatAction {
	return &eatAction{
		Action: NewAction("eat", cost),
	}
}

type eatAction struct {
	Action
	inRange bool
}

func (a *eatAction) Reset() {}

func (a *eatAction) CheckProceduralPrecondition(agent Agent) bool {
	return true
}

func (a *eatAction) RequiresInRange() bool {
	return false
}

func (a *eatAction) IsInRange() bool {
	return true
}

func (a *eatAction) Perform(agent Agent) bool {
	return true
}

func (a *eatAction) IsDone() bool {
	return true
}

func newSleepAction(cost float64) *sleepAction {
	return &sleepAction{
		Action: NewAction("sleep", cost),
	}
}

type sleepAction struct {
	Action
}

func (a *sleepAction) Reset() {}

func (a *sleepAction) CheckProceduralPrecondition(agent Agent) bool {
	return true
}

func (a *sleepAction) RequiresInRange() bool {
	return false
}

func (a *sleepAction) IsInRange() bool {
	return true
}

func (a *sleepAction) Perform(agent Agent) bool {
	return true
}

// dataProvider interfaces with the world
type dataProvider struct {
	moveResult bool
}

func (p *dataProvider) GetWorldState() StateList {
	worldState := make(StateList, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false
	return worldState
}

func (p *dataProvider) CreateGoalState() StateList {
	goal := make(StateList, 0)
	goal["isFull"] = true
	return goal
}

func (p *dataProvider) PlanFailed(failedGoal StateList) {
	fmt.Printf("Planning failed: %v\n", failedGoal)
}

func (p *dataProvider) PlanFound(goal StateList, actions []Actionable) {
	fmt.Printf("Planning success to goal %v with actions %v\n", goal, actions)
}

func (p *dataProvider) ActionsFinished() {
	fmt.Println("all actions finished")
}

func (p *dataProvider) PlanAborted(aborter Actionable) {
	fmt.Printf("plan aborted by %v\n", aborter)
}

func (p *dataProvider) MoveAgent(nextAction Actionable) bool {
	return p.moveResult
}
