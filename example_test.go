package goap

import "fmt"

var (
	HaveFood = State{"hasFood", true}
	Hungry   = State{"isFull", false}
	Tired    = State{"isTired", true}
	IsHurt   = State{"isHurt", true}
	IsHidden = State{"isHidden", true}
)

func findFood() *testAction {
	findFood := newTestAction("getFood", 8, false)
	findFood.AddEffect(HaveFood)
	findFood.AddPrecondition(Dont(HaveFood))
	return findFood
}

func eatAction() *testAction {
	eat := newTestAction("eat", 4, false)
	eat.AddEffect(Isnt(Hungry), Dont(HaveFood))
	eat.AddPrecondition(HaveFood, Hungry)
	return eat
}

func sleepAction() *testAction {
	sleep := newTestAction("sleep", 4, false)
	sleep.AddEffect(Isnt(Tired))
	sleep.AddPrecondition(Tired)
	return sleep
}

type TestAgent struct {
	DefaultAgent
	frame int
}

func (p *TestAgent) PlanFailed(failedGoal StateList) {
	fmt.Printf("Planning failed: %v\n", failedGoal)
}

func (p *TestAgent) PlanFound(goal StateList, actions []Action) {
	fmt.Printf("Planning success to goal %v with actions %v\n", goal, actions)
}

func (p *TestAgent) ActionsFinished() {
	fmt.Println("all actions finished")
}

func (p *TestAgent) PlanAborted(aborter Action) {
	fmt.Printf("plan aborted by %v\n", aborter)
}

func (p *TestAgent) MoveAgent(nextAction Action) bool {
	return nextAction.IsInRange()
}

func (a *TestAgent) Update() {
	a.frame++
	fmt.Printf("#%d\n", a.frame)
	a.FSM(a, func(m string) {
		fmt.Println(m)
	})
}

func ExamplePlan() {
	getFood := newGetFoodAction(8)
	getFood.AddEffect(HaveFood)

	eat := newEatAction(4)
	eat.AddPrecondition(HaveFood, Hungry)
	eat.AddEffect(Isnt(Hungry))

	sleep := newSleepAction(4)
	sleep.AddPrecondition(Tired)
	sleep.AddEffect(Isnt(Tired))

	agent := TestAgent{
		DefaultAgent: NewDefaultAgent([]Action{getFood, eat, sleep}),
	}

	currentState := make(StateList)
	currentState.Add(Hungry).Dont(HaveFood)
	agent.SetState(currentState)

	goal := make(StateList)
	goal.Isnt(Hungry)
	agent.SetGoalState(goal)

	// 1. idle state, will do planning, plan will be to getFood and then eat it
	agent.Update()

	// 2. Discovers that it needs to move to get food
	agent.Update()

	// 3. Move to food, it instantly succeeds
	getFood.inRange = true
	agent.Update()

	// 4. We have moved and food is in range, so run getFood action
	agent.Update()

	// 5. getFood action as done, progress to eat action
	agent.Update()

	// 6. eat action is done, there are no more steps in the plan
	agent.Update()

	// Output:
	//#1
	// Idle - is planning
	// Planning success to goal map[isFull:true] with actions [getFood eat]
	// #2
	// Do - scheduling moveTo getFood
	// #3
	// MoveTo - MoveAgent(getFood)
	// MoveTo - done
	// #4
	// Do - getFood.Perform()
	// #5
	// Do - action getFood is done
	// Do - eat.Perform()
	// #6
	// Do - action eat is done
	// all actions finished
}
