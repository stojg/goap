package goap

func NewDefaultAgent(actions []Actionable) DefaultAgent {
	return DefaultAgent{
		availableActions: actions,
		StateMachine:     NewFSM(Idle),
	}
}

type DefaultAgent struct {
	StateMachine *FSM

	availableActions []Actionable
	currentActions   []Actionable

	currentState StateList
	goalState    StateList

	moveResult bool
}

func (a *DefaultAgent) AvailableActions() []Actionable {
	return a.availableActions
}

func (a *DefaultAgent) AddAction(action Actionable) {
	a.availableActions = append(a.availableActions, action)
}

func (a *DefaultAgent) CurrentActions() []Actionable {
	return a.currentActions
}

func (a *DefaultAgent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

func (a *DefaultAgent) SetCurrentActions(actions []Actionable) {
	a.currentActions = actions
}

func (a *DefaultAgent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

func (a *DefaultAgent) State() StateList {
	return a.currentState
}

func (a *DefaultAgent) SetState(l StateList) {
	a.currentState = l
}

func (a *DefaultAgent) GoalState() StateList {
	return a.goalState
}

func (a *DefaultAgent) SetGoalState(l StateList) {
	a.goalState = l
}

func (p *DefaultAgent) PlanFailed(failedGoal StateList) {}

func (p *DefaultAgent) PlanFound(goal StateList, actions []Actionable) {}

func (p *DefaultAgent) ActionsFinished() {}

func (p *DefaultAgent) PlanAborted(aborter Actionable) {}

func (p *DefaultAgent) MoveAgent(nextAction Actionable) bool {
	return false
}

func (a *DefaultAgent) FSM(b Agent, debug func(string)) {
	a.StateMachine.Update(b, debug)
}

func (a *DefaultAgent) Update() {
	a.FSM(a, func(m string) {})
}
