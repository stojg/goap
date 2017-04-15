package goap

func newGetFoodAction(cost float64) *getFoodAction {
	a := &getFoodAction{
		Action: NewInRangeAction("getFood", cost),
	}
	return a
}

type getFoodAction struct {
	Action
	hasFood bool
}

func (a *getFoodAction) CheckContextPrecondition(agent Agent) bool {
	a.SetTarget([]int{10, 0, 200})
	return true
}

func (a *getFoodAction) Perform(agent Agent) bool {
	a.hasFood = true
	return true
}

func (a *getFoodAction) IsDone() bool {
	return a.hasFood
}

func newEatAction(cost float64) *eatingAction {
	a := &eatingAction{
		Action: NewAction("eat", cost),
	}
	return a
}

type eatingAction struct {
	Action
}

func (a *eatingAction) Perform(agent Agent) bool {
	return true
}

func (a *eatingAction) IsDone() bool {
	return true
}

func newSleepAction(cost float64) *sleepingAction {
	return &sleepingAction{
		Action: NewAction("sleep", cost),
	}
}

type sleepingAction struct {
	Action
}

func (a *sleepingAction) IsInRange() bool {
	return true
}

func (a *sleepingAction) Perform(agent Agent) bool {
	return true
}
