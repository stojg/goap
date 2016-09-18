package goap

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
