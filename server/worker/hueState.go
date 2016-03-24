package worker

import (
    "github.com/eugene-kartsev/go-teamcity-hue/server/hue"
    "fmt"
)

type hueState struct {
    id string
    tcStates map[string]bool
}

func (state hueState) changeState(tcId string, newState bool, canReed chan bool) {
    <-canReed

    if currentTcState, found := state.tcStates[tcId]; found {
        if currentTcState != newState {
            allOk := true
            for _, isStateOk := range state.tcStates {
                if !isStateOk {
                    allOk = false
                }
            }
            if hueNode, found := hueList[state.id]; found {
                var signal int
                if allOk {
                    signal = hue.GREEN
                } else {
                    signal = hue.RED
                }
                _, err := hue.Signal(hueNode.Url, signal)
                if err != nil {
                    fmt.Println("Cannot process hue request ["+hueNode.Url+"]")
                } else {
                    state.tcStates[tcId] = newState;
                }
            }
        }
    }

    canReed<-true
}
