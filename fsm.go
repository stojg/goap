package goap

import (
	"fmt"
)

type FSMState func(fsm *FSM, obj Agent, debug func(string))

func NewFSM(startState FSMState) *FSM {
	fsm := &FSM{}
	fsm.Reset(startState)
	return fsm
}

type FSM struct {
	stateStack []FSMState
}

func (fsm *FSM) Update(agent Agent, debug func(string)) {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack[len(fsm.stateStack)-1](fsm, agent, debug)
	}
}

func (fsm *FSM) Push(state FSMState) {
	fsm.stateStack = append(fsm.stateStack, state)
}

func (fsm *FSM) Reset(state FSMState) {
	states := len(fsm.stateStack)
	for i := 0; i < states; i++ {
		fsm.Pop()
	}
	fsm.Push(state)
}

func (fsm *FSM) Pop() {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack = fsm.stateStack[:len(fsm.stateStack)-1]
	}
}

func Idle(fsm *FSM, agent Agent, debug func(string)) {
	debug("Idle - is planning")
	goal := agent.GoalState()
	plan := Plan(agent, agent.AvailableActions(), agent.State(), goal)
	if plan == nil {
		agent.PlanFailed(goal)
		return
	}
	agent.SetCurrentActions(plan)
	agent.PlanFound(goal, plan)
	fsm.Reset(Do)
}

func Do(fsm *FSM, agent Agent, debug func(string)) {
	// no actions to perform
	if len(agent.CurrentActions()) == 0 {
		fsm.Reset(Idle)
		agent.ActionsFinished()
		return
	}

	action := agent.CurrentActions()[0]
	if action.IsDone() {
		debug(fmt.Sprintf("Do - action %s is done", action))
		// the action is done. Remove it so we can perform the next one
		agent.PopCurrentAction()
	}

	if len(agent.CurrentActions()) == 0 {
		fsm.Reset(Idle)
		agent.ActionsFinished()
		return
	}

	action = agent.CurrentActions()[0]
	inRange := true
	if action.RequiresInRange() {
		inRange = action.IsInRange()
	}
	// we need to move there first
	if !inRange {
		debug(fmt.Sprintf("Do - scheduling moveTo %s", action))
		fsm.Push(MoveTo)
		return
	}

	// we are in range, so perform the action
	debug(fmt.Sprintf("Do - %s.Perform()", action))
	success := action.Perform(agent)
	// action failed, we need to plan again
	if !success {
		fsm.Reset(Idle)
		agent.PlanAborted(action)
	}
}

func MoveTo(fsm *FSM, agent Agent, debug func(string)) {
	action := agent.CurrentActions()[0]

	if action.Target() == nil {
		debug("Error: MoveTo requires a target but has none. Planning failed. You did not assign the target in your Action.CheckContextPrecondition()")
		fsm.Reset(Idle)
		return
	}

	// get the agent to move itself
	debug(fmt.Sprintf("MoveTo - MoveAgent(%s)", action))
	if agent.MoveAgent(action) {
		debug("MoveTo - done")
		fsm.Pop()
	}
}
