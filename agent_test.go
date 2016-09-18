package goap

import "fmt"

func newTestAgent(actions []Actionable) *testAgent {
	agent := &testAgent{
		availableActions: actions,
		fsm:              &FSM{},
	}

	agent.idleState = func(fsm *FSM, obj Agent) {
		worldState := agent.GetWorldState()
		goal := agent.CreateGoalState()
		agent.debugf("idle - is planning\n")
		plan := Plan(obj, agent.availableActions, worldState, goal)
		if plan != nil {
			agent.SetCurrentActions(plan)
			agent.PlanFound(goal, plan)
			// move to PerformAction state
			fsm.PopState()
			fsm.PushState(agent.doAction)
		} else {
			agent.PlanFailed(goal)
			// move back to IdleAction state
			fsm.PopState()
			fsm.PushState(agent.idleState)
		}
	}

	agent.moveToState = func(fsm *FSM, obj Agent) {
		action := agent.CurrentAction()
		if action.RequiresInRange() && action.Target() == nil {
			agent.debugf("Error: Action requires a target but has none. Planning failed. You did not assign the target in your Action.CheckProceduralPrecondition()\n")
			fsm.PopState() // move
			fsm.PopState() // perform
			fsm.PushState(agent.idleState)
			return
		}

		// get the agent to move itself
		agent.debugf("moveTo - MoveAgent(%s)\n", action)
		if agent.MoveAgent(action) {
			agent.debugf("moveTo - done\n")
			fsm.PopState()
		}
	}

	agent.doAction = func(fsm *FSM, obj Agent) {
		// no actions to perform
		if !agent.HasActionPlan() {
			fsm.PopState()
			fsm.PushState(agent.idleState)
			agent.ActionsFinished()
			return
		}

		action := agent.CurrentAction()
		if action.IsDone() {
			agent.debugf("doAction - action %s is done\n", action)
			// the action is done. Remove it so we can perform the next one
			agent.PopCurrentAction()
		}

		if agent.HasActionPlan() {
			// perform the next action
			action = agent.CurrentAction()
			inRange := true
			if action.RequiresInRange() {
				inRange = action.IsInRange()
			}
			if inRange {
				// we are in range, so perform the action
				agent.debugf("doAction - %s.Perform()\n", action)
				success := action.Perform(obj)
				if !success {
					// action failed, we need to plan again
					fsm.PopState()
					fsm.PushState(agent.idleState)
					agent.PlanAborted(action)
				}
			} else {
				agent.debugf("doAction - scheduling moveTo %s\n", action)
				// we need to move there first
				fsm.PushState(agent.moveToState)
			}

		} else {
			// no actions left, move to Plan state
			fsm.PopState()
			fsm.PushState(agent.idleState)
			agent.ActionsFinished()
		}
	}
	agent.fsm.PushState(agent.idleState)
	return agent
}

type testAgent struct {
	Debug bool

	fsm   *FSM
	frame int

	idleState   FSMState
	moveToState FSMState
	doAction    FSMState

	availableActions []Actionable
	currentActions   []Actionable

	moveResult bool
}

func (a *testAgent) StateMachine() *FSM {
	return a.fsm
}

func (a *testAgent) AvailableActions() []Actionable {
	return a.availableActions
}

func (a *testAgent) AddAction(action Actionable) {
	a.availableActions = append(a.availableActions, action)
}

func (a *testAgent) CurrentAction() Actionable {
	return a.currentActions[0]
}

func (a *testAgent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

func (a *testAgent) SetCurrentActions(actions []Actionable) {
	a.currentActions = actions
}

func (a *testAgent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

func (a *testAgent) debugf(format string, v ...interface{}) {
	if a.Debug {
		fmt.Printf(format, v...)
	}
}

func (p *testAgent) GetWorldState() StateList {
	worldState := make(StateList, 0)
	worldState["isFull"] = false
	worldState["hasFood"] = false
	return worldState
}

func (p *testAgent) CreateGoalState() StateList {
	goal := make(StateList, 0)
	goal["isFull"] = true
	return goal
}

func (p *testAgent) PlanFailed(failedGoal StateList) {
	fmt.Printf("Planning failed: %v\n", failedGoal)
}

func (p *testAgent) PlanFound(goal StateList, actions []Actionable) {
	fmt.Printf("Planning success to goal %v with actions %v\n", goal, actions)
}

func (p *testAgent) ActionsFinished() {
	fmt.Println("all actions finished")
}

func (p *testAgent) PlanAborted(aborter Actionable) {
	fmt.Printf("plan aborted by %v\n", aborter)
}

func (p *testAgent) MoveAgent(nextAction Actionable) bool {
	return p.moveResult
}

func (a *testAgent) Update() {
	a.frame++
	if a.Debug {
		fmt.Printf("#%d\n", a.frame)
	}
	a.StateMachine().Update(a)
}
