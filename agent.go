package goap

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
	planFailed(failedGoal StateList)

	/**
	 * A plan was found for the supplied goal.
	 * These are the actions the Agent will perform, in order.
	 */
	planFound(goal StateList, actions []Actionable)

	/**
	 * All actions are complete and the goal was reached. Hooray!
	 */
	actionsFinished()

	/**
	 * One of the actions caused the plan to abort.
	 * That action is returned.
	 */
	planAborted(aborter Actionable)

	/**
	 * Called during Update. Move the agent towards the target in order
	 * for the next action to be able to perform.
	 * Return true if the Agent is at the target and the next action can perform.
	 * False if it is not there yet.
	 */
	moveAgent(nextAction Actionable) bool
}
