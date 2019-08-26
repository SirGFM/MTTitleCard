package main

import (
    "github.com/SirGFM/MTTitleCard/config"
    "github.com/SirGFM/MTTitleCard/page"
    "flag"
    "fmt"
    "os"
    "os/signal"
)

func main() {
    var configFile string

    flag.StringVar(&configFile, "config", "", "Path to a config file")
    flag.Parse()
    err := config.Load(configFile)
    if err != nil {
        panic(fmt.Sprintf("Failed to load the configuratin: %+v", err))
    }

    err = page.StartServer(config.Get().Port)
    if err != nil {
        panic(fmt.Sprintf("Failed to start server: %+v", err))
    }

    signalTrap := make(chan os.Signal, 1)
    wait := make(chan struct{}, 1)
    go func (c chan os.Signal) {
        _ = <-c
        page.StopServer()
        wait <- struct{}{}
    } (signalTrap)
    signal.Notify(signalTrap, os.Interrupt)

    <-wait
}
