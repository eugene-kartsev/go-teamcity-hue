package worker

type tcState struct {
    id string
    currentState bool
    hueList []string
}

func (state tcState) changeState(newState bool, canReed chan bool) {
    <-canReed

    if newState != state.currentState {
        for _, hueId := range state.hueList {
            if hueState, found := hueStates[hueId]; found {
                go hueState.changeState(state.id, newState, canReed)
            }
        }
        state.currentState = newState
    }

    canReed<-true
}

