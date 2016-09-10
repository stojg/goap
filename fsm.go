package goap

type FSMState func(fsm *FSM, obj GameObject)

type FSM struct {
	stateStack []FSMState
}

func (fsm *FSM) Update(gameObject interface{}) {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack[len(fsm.stateStack)-1](fsm, gameObject)
	}
}

func (fsm *FSM) PushState(state FSMState) {
	fsm.stateStack = append(fsm.stateStack, state)
}

func (fsm *FSM) PopState() {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack = fsm.stateStack[:len(fsm.stateStack)-1]
	}
}
