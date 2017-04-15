package goap

func NewDefaultAgent(actions []Action) DefaultAgent {
	return DefaultAgent{
		availableActions: actions,
		StateMachine:     NewFSM(Idle),
	}
}

type DefaultAgent struct {
	StateMachine *FSM

	availableActions []Action
	currentActions   []Action

	currentState StateList
	goalState    StateList

	moveResult bool
}

func (a *DefaultAgent) AvailableActions() []Action {
	return a.availableActions
}

func (a *DefaultAgent) AddAction(action Action) {
	a.availableActions = append(a.availableActions, action)
}

func (a *DefaultAgent) CurrentActions() []Action {
	return a.currentActions
}

func (a *DefaultAgent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

func (a *DefaultAgent) SetCurrentActions(actions []Action) {
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

func (a *DefaultAgent) AddState(s State) {
	a.currentState[s.Name] = s.Value
}

func (a *DefaultAgent) GoalState() StateList {
	return a.goalState
}

func (a *DefaultAgent) SetGoalState(l StateList) {
	a.goalState = l
}

func (p *DefaultAgent) PlanFailed(failedGoal StateList) {}

func (p *DefaultAgent) PlanFound(goal StateList, actions []Action) {}

func (p *DefaultAgent) ActionsFinished() {}

func (p *DefaultAgent) PlanAborted(aborter Action) {}

func (p *DefaultAgent) MoveAgent(nextAction Action) bool {
	return false
}

func (a *DefaultAgent) FSM(b Agent, debug func(string)) {
	a.StateMachine.Update(b, debug)
}

func (a *DefaultAgent) Update() {
	a.FSM(a, func(m string) {})
}
