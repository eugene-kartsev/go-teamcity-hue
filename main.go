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

    worker.Init(cfg);

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