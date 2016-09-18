package goap

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
	agent := newTestAgent(actions)
	agent.Debug = true

	// 1. idle state, will do planning
	agent.Update()

	// 2. perform action getFood, but discovers that will need to move
	agent.Update()

	// 3. Move to food, it' instantly succeeds
	agent.Update()

	// 4. We have moved and food is in range
	agent.moveResult = true
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
