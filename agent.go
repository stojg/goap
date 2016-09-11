package goap

import (
	"fmt"
)

func NewGoapAgent(dataProvider DataProvider, actions []Actionable) *GoapAgent {
	agent := &GoapAgent{
		fsm:              &FSM{},
		dataProvider:     dataProvider,
		availableActions: actions,
	}

	agent.idleState = func(fsm *FSM, obj Agent) {
		worldState := agent.dataProvider.GetWorldState()
		goal := agent.dataProvider.CreateGoalState()
		agent.debugf("idle - is planning\n")
		plan := Plan(obj, agent.availableActions, worldState, goal)
		if plan != nil {
			agent.SetCurrentActions(plan)
			agent.dataProvider.PlanFound(goal, plan)
			// move to PerformAction state
			fsm.PopState()
			fsm.PushState(agent.doAction)
		} else {
			agent.dataProvider.PlanFailed(goal)
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
		if agent.dataProvider.MoveAgent(action) {
			agent.debugf("moveTo - done\n")
			fsm.PopState()
		}
	}

	agent.doAction = func(fsm *FSM, obj Agent) {
		// no actions to perform
		if !agent.HasActionPlan() {
			fsm.PopState()
			fsm.PushState(agent.idleState)
			agent.dataProvider.ActionsFinished()
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
					agent.dataProvider.PlanAborted(action)
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
			agent.dataProvider.ActionsFinished()
		}
	}
	agent.fsm.PushState(agent.idleState)
	return agent
}

type GoapAgent struct {
	Debug bool

	fsm   *FSM
	frame int

	idleState   FSMState
	moveToState FSMState
	doAction    FSMState

	availableActions []Actionable
	currentActions   []Actionable

	dataProvider DataProvider
}

func (a *GoapAgent) StateMachine() *FSM {
	return a.fsm
}

func (a *GoapAgent) Update() {
	a.fsm.Update(a)
}

func (a *GoapAgent) AddAction(action Actionable) {
	a.availableActions = append(a.availableActions, action)
}

func (a *GoapAgent) CurrentAction() Actionable {
	return a.currentActions[0]
}

func (a *GoapAgent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

func (a *GoapAgent) SetCurrentActions(actions []Actionable) {
	a.currentActions = actions
}

func (a *GoapAgent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

func (a *GoapAgent) debugf(format string, v ...interface{}) {
	if a.Debug {
		fmt.Printf(format, v...)
	}
}

type Agent interface{}

type GameObject interface{}

// IGoap must be implemented by agents that wants to use GOAP  It provides information to the GOAP
// planner so it can plan what actions to use.
//
// It also provides an interface for the planner to give feedback to the Agent and report
// success/failure.
type DataProvider interface {
	/**
	 * The starting state of the Agent and the world.
	 * Supply what states are needed for actions to run.
	 */
	GetWorldState() StateList

	/**
	 * Give the planner a new goal so it can figure out
	 * the actions needed to fulfill it.
	 */
	CreateGoalState() StateList

	/**
	 * No sequence of actions could be found for the supplied goal.
	 * You will need to try another goal
	 */
	PlanFailed(failedGoal StateList)

	/**
	 * A plan was found for the supplied goal.
	 * These are the actions the Agent will perform, in order.
	 */
	PlanFound(goal StateList, actions []Actionable)

	/**
	 * All actions are complete and the goal was reached. Hooray!
	 */
	ActionsFinished()

	/**
	 * One of the actions caused the plan to abort.
	 * That action is returned.
	 */
	PlanAborted(aborter Actionable)

	/**
	 * Called during Update. Move the agent towards the target in order
	 * for the next action to be able to perform.
	 * Return true if the Agent is at the target and the next action can perform.
	 * False if it is not there yet.
	 */
	MoveAgent(nextAction Actionable) bool
}
