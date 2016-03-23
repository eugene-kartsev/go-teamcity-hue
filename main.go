package main

import (
    "github.com/eugene-kartsev/go-teamcity-hue/server/worker"
    "fmt"
)

func main() {
    tcUrl := ""
    tcLogin := ""
    tcPassword := ""
    hueUrl := ""
    interval := 60

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