package main

import (
    "fmt"
    "go-teamcity-hue/teamcity"
    "go-teamcity-hue/hue"
    "go-teamcity-hue/dispatcher"
    "go-teamcity-hue/config"
)

func main() {

    cfg, err := config.Create()
    if err != nil {
        return
    }


    tc, err := teamcity.Create(cfg)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    hue, err := hue.Create(cfg)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println("CONFIGURATION is OK. Starting...")
    go dispatcher.Start(tc, hue)

    waitEnterThenQuit()
}

func waitEnterThenQuit() {
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