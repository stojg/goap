package goap

// Agent must be implemented by agents that wants to use GOAP  It provides information to the GOAP
// planner so it can plan what actions to use.
//
// It also provides an interface for the planner to give feedback to the Agent and report
// success/failure.
type Agent interface {

	// Get the actions that this agent can do
	AvailableActions() []Actionable

	// Set the actions that will allow this agent to reach it's goal
	SetCurrentActions([]Actionable)

	// Get current planned actions
	CurrentActions() []Actionable

	// Move the agent towards the target in order for the next action to be able to perform. Return true if the Agent
	// is at the target and the next action can perform, false if it is not there yet.
	MoveAgent(Actionable) bool

	// Advance the internal state machine and run actions
	Update()

	// Remove the currently running Action
	PopCurrentAction()

	// The starting state of the Agent and the world. Supplies what states are needed for actions to run.
	State() StateList

	// Set the starting state of the Agent, includes world information.
	SetState(StateList)

	AddState(State)

	// Get the goal for this actor
	GoalState() StateList

	// Set the goal for this actor
	SetGoalState(StateList)

	// Below are life-cycle hooks that are called during the different stages
	// of the planning.

	// No sequence of actions could be found for the supplied goal. You will need to try another goal
	PlanFailed(failedGoal StateList)

	// A plan was found for the supplied goal. These are the actions the Agent will perform, in order.
	PlanFound(goal StateList, actions []Actionable)

	// All actions are complete and the goal was reached.
	ActionsFinished()

	// One of the actions caused the plan to abort. That action is passed in.
	PlanAborted(Actionable)
}
