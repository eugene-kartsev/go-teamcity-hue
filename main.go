package main

import (
    "github.com/eugene-kartsev/go-teamcity-hue/server/worker"
    "github.com/eugene-kartsev/go-teamcity-hue/server/config"
    "fmt"
)

func main() {

    cfg, err := config.Read()
    if err != nil {
        fmt.Println(err)
        return
    }

    tcUrl := cfg.TeamCityNodes[0].Url
    tcLogin := cfg.TeamCityNodes[0].Login
    tcPassword := cfg.TeamCityNodes[0].Password
    hueUrl := cfg.HueNodes[0].Url
    interval := cfg.TeamCityNodes[0].Interval

    go worker.Start(tcUrl, tcLogin, tcPassword, hueUrl, interval)

    onQuit()
}

func onQuit() {
    quit := make(chan bool)

    go func() {
        fmt.Println("Press ENTER to close the app...")

        var message string
        fmt.Scanf("%s", &message)
        fmt.Println("I'm done here. Closing.")
        quit<-true
    }()

    <-quit
}