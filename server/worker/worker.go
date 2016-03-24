package worker

import (
    "github.com/eugene-kartsev/go-teamcity-hue/server/teamcity"
    "github.com/eugene-kartsev/go-teamcity-hue/server/config"
    "time"
    "fmt"
)

const (
    isTeamcitySuccessByDefault = false
)

var hueList = make(map[string]config.HueNode)
var tcList  = make(map[string]config.TeamCityNode)

var tcStates  = make(map[string]tcState)
var hueStates = make(map[string]hueState)

func Init(cfg *config.Config) {

    var canReed = make(chan bool, 1)
    canReed <- true

    for _, hueNode := range cfg.HueNodes {
        if _, found := hueList[hueNode.Id]; !found {
            hueList[hueNode.Id] = hueNode
            hueStates[hueNode.Id] = hueState{
                id:hueNode.Id,
                tcStates:make(map[string]bool),
            }
        }
    }

    for _, tcNode := range cfg.TeamCityNodes {
        if _, found := tcList[tcNode.Id]; !found {
            tcList[tcNode.Id] = tcNode
            tcStates[tcNode.Id] = tcState{
                id:tcNode.Id,
                currentState:isTeamcitySuccessByDefault,
                hueList:make([]string, 0),
            }
        }
    }

    for _, mapNode := range cfg.Map {
        hueId := mapNode.HueId
        if _, hueFound := hueList[hueId]; hueFound {
            for _, tcId := range mapNode.TcIds {
                if tc, tcFound := tcList[tcId]; tcFound {
                    hueStates[hueId].tcStates[tcId] = isTeamcitySuccessByDefault

                    state := tcStates[tcId]
                    state.hueList =  append(tcStates[tcId].hueList, hueId)

                    go startWorker(tc.Url, tc.Login, tc.Password, tc.Interval, &state, canReed);
                }
            }
        }
    }
}

func startWorker(tcUrl string, tcLogin string, tcPassword string, interval int, state *tcState, canReed chan bool) {
    go func() {
        for {
            tcStatus, tcErr := teamcity.GetBuildStatus(tcUrl, tcLogin, tcPassword)
            if tcErr != nil {
                fmt.Println("Cannot connect to TeamCity ["+tcUrl+"]. Error: " + tcErr.Error())
            } else {
                fmt.Println("TeamCity ["+tcUrl+"] STATUS: " + tcStatus)
                go state.changeState(tcStatus == teamcity.SUCCESS, canReed)
            }

            time.Sleep(time.Duration(interval) * time.Second)
        }
    }()
}


