package goap

import (
	"fmt"
)

func NewExampleAgent(dataProvider DataProvider, actions []Actionable) *ExampleAgent {

	a := &ExampleAgent{
		dataProvider:     dataProvider,
		availableActions: actions,
	}

	a.idleState = func(fsm *FSM, obj GameObject) {
		fmt.Println("idleState")
		worldState := a.dataProvider.GetWorldState()
		goal := a.dataProvider.CreateGoalState()
		plan := Plan(obj, a.availableActions, worldState, goal)
		if plan != nil {
			a.currentActions = plan
			a.dataProvider.planFound(goal, plan)
			// move to PerformAction state
			fsm.PopState()
			fsm.PushState(a.performActionState)
		} else {
			a.dataProvider.planFailed(goal)
			// move back to IdleAction state
			fsm.PopState()
			fsm.PushState(a.idleState)
		}
	}

	a.moveToState = func(fsm *FSM, obj GameObject) {
		action := a.currentActions[0]
		fmt.Printf("moveToState %s\n", action)
		if action.RequiresInRange() && action.Target() == nil {
			fmt.Println("Error: Action requires a target but has none. Planning failed. You did not assign the target in your Action.CheckProceduralPrecondition()")
			fsm.PopState() // move
			fsm.PopState() // perform
			fsm.PushState(a.idleState)
			return
		}

		// get the agent to move itself
		if a.dataProvider.moveAgent(action) {
			fmt.Printf("done moveToState %s\n", action)
			fsm.PopState()
		}
	}

	a.performActionState = func(fsm *FSM, obj GameObject) {
		// no actions to perform
		if !a.HasActionPlan() {
			fsm.PopState()
			fsm.PushState(a.idleState)
			a.dataProvider.actionsFinished()
			return
		}

		action := a.currentActions[0]
		fmt.Printf("performActionState %s\n", action)
		if action.IsDone() {
			// the action is done. Remove it so we can perform the next one
			a.currentActions = a.currentActions[1:]
		}

		if a.HasActionPlan() {
			// perform the next action
			action = a.currentActions[0]
			inRange := true

			if action.RequiresInRange() {
				inRange = action.IsInRange()
			}
			if inRange {
				// we are in range, so perform the action
				success := action.Perform(obj)

				if !success {
					// action failed, we need to plan again
					fsm.PopState()
					fsm.PushState(a.idleState)
					a.dataProvider.planAborted(action)
				}
			} else {

				// we need to move there first
				// push moveTo state
				fsm.PushState(a.moveToState)
			}

		} else {
			// no actions left, move to Plan state
			fsm.PopState()
			fsm.PushState(a.idleState)
			a.dataProvider.actionsFinished()
		}
	}

	a.fsm.PushState(a.idleState)

	return a
}

// ExampleAgent inspired by
// https://github.com/sploreg/goap/blob/d1cea0728fb4733266affea8049da1e373d618f7/Assets/Standard%20Assets/Scripts/AI/GOAP/GoapAgent.cs
type ExampleAgent struct {
	fsm   FSM
	frame int

	idleState          FSMState
	moveToState        FSMState
	performActionState FSMState

	availableActions []Actionable
	currentActions   []Actionable
	// this is the implementing class that provides our world data and listens to feedback on planning
	dataProvider DataProvider
}

func (a *ExampleAgent) Update() {
	a.frame++
	fmt.Printf("Update #%d\n", a.frame)
	a.fsm.Update(a)

}

func (a *ExampleAgent) AddAction(action Actionable) {
	a.availableActions = append(a.availableActions, action)
}

func (a *ExampleAgent) getAction(action Actionable) {
}

func (a *ExampleAgent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

type TestDataProvider struct {
	moveResult bool
}

func (p *TestDataProvider) GetWorldState() StateList {
	worldState := make(StateList, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false
	return worldState
}

func (p *TestDataProvider) CreateGoalState() StateList {
	goal := make(StateList, 0)
	goal["isFull"] = true
	return goal
}

func (p *TestDataProvider) planFailed(failedGoal StateList) {
	fmt.Printf("Planning failed: %v\n", failedGoal)
}

func (p *TestDataProvider) planFound(goal StateList, actions []Actionable) {
	fmt.Printf("Planning success to goal %v with actions %v\n", goal, actions)
}

func (p *TestDataProvider) actionsFinished() {
	fmt.Println("actions finished")
}

func (p *TestDataProvider) planAborted(aborter Actionable) {
	fmt.Printf("plan aborted by %v\n", aborter)
}

func (p *TestDataProvider) moveAgent(nextAction Actionable) bool {
	return p.moveResult
}

func ExamplePlan() {
	getFood := newTestAction("getFood", 8, true)
	getFood.AddPrecondition("hasFood", false)
	getFood.isDone = false
	getFood.target = []int{1, 0, 1}
	getFood.AddEffect("hasFood", true)

	eat := newTestAction("eat", 4, false)
	eat.AddPrecondition("hasFood", true)
	eat.AddPrecondition("isFull", false)
	eat.AddEffect("isFull", true)

	sleep := newTestAction("sleep", 4, false)
	sleep.AddPrecondition("isTired", true)
	sleep.AddEffect("isTired", false)

	actions := []Actionable{getFood, eat, sleep}
	provider := &TestDataProvider{}
	agent := NewExampleAgent(provider, actions)

	// 1. idle state, will do planning
	agent.Update()

	// 2. perform action getFood, but discovers that will need to move
	agent.Update()

	// 3. Move to food, it' instantly succeeds
	provider.moveResult = true
	getFood.inRange = true
	agent.Update()

	// 4. mark the getFoodAction as done
	getFood.isDone = true
	agent.Update()

	// 5. time to eat that food
	agent.Update()

	// Output:
	// Update #1
	// idleState
	// Planning success to goal map[isFull:true] with actions [getFood eat]
	// Update #2
	// performActionState getFood
	// Update #3
	// moveToState getFood
	// done moveToState getFood
	// Update #4
	// performActionState getFood
	// Update #5
	// performActionState eat
	// actions finished

}
