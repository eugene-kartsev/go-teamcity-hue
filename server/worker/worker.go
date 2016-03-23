package worker

import (
    "github.com/eugene-kartsev/go-teamcity-hue/server/teamcity"
    "github.com/eugene-kartsev/go-teamcity-hue/server/hue"
    "time"
    "fmt"
)

func Start(tcUrl string, tcLogin string, tcPassword string, hueUrl string, interval int)  {
    go func() {
        for {
            tcStatus, tcErr := teamcity.GetBuildStatus(tcUrl, tcLogin, tcPassword)
            if tcErr != nil {
                fmt.Println("TeamCity error: " + tcErr.Error())
            } else {
                fmt.Println("TeamCity status: " + tcStatus)
                if (tcStatus == teamcity.SUCCESS) {
                    hueStatus, hueErr := hue.Signal(hueUrl, hue.GREEN)
                    if !hueStatus {
                        fmt.Println("HUE error: " + hueErr.Error())
                    }
                } else {
                    hueStatus, hueErr := hue.Signal(hueUrl, hue.RED)
                    if !hueStatus {
                        fmt.Println("HUE error: " + hueErr.Error())
                    }
                }
            }

            time.Sleep(time.Duration(interval) * time.Second)
        }
    }()
}


