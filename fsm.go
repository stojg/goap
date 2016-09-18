package goap

type FSMState func(fsm *FSM, obj Agent)

type FSM struct {
	stateStack []FSMState
}

func (fsm *FSM) Update(agent Agent) {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack[len(fsm.stateStack)-1](fsm, agent)
	}
}

func (fsm *FSM) PushState(state FSMState) {
	fsm.stateStack = append(fsm.stateStack, state)
}

func (fsm *FSM) Clear() {
	states := len(fsm.stateStack)
	for i := 0; i < states; i++ {
		fsm.PopState()
	}
}

func (fsm *FSM) PopState() {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack = fsm.stateStack[:len(fsm.stateStack)-1]
	}
}
