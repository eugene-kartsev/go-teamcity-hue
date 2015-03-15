package main

import (
    "fmt"
    "go-teamcity-hue/teamcity"
    "go-teamcity-hue/hue"
    "time"
)

func main() {
    config := teamcity.Config{}
    config.Login    = "ca\\yevgenk"
    config.Password = "5tgb%TGB6yhn^YHN"
    config.RefreshInSec = 20 * time.Second

    api := teamcity.Setup(config)

    hue.Init("http://10.10.0.80/api/newdeveloper/lights/")

    api.Watch(func(statusOK bool) {
        if statusOK {
            //fmt.Println("GREEN")
            go hue.Green()
        } else {
            //fmt.Println("RED")
            go hue.Red()
        }
    })

    quit := make(chan bool)

    go func() {
        fmt.Println("Press Enter to close the app...")

        var message string
        fmt.Scanf("%s", &message)
        fmt.Println("I'm done here. Closing.")
        quit<-true
    }()

    <-quit
}
