package main

import (
    "github.com/SirGFM/MTTitleCard/config"
    "github.com/SirGFM/MTTitleCard/page"
    "flag"
    "log"
    "log/syslog"
    "os"
    "os/signal"
)

func main() {
    var configFile string

    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)

    sl, err := syslog.New(syslog.LOG_LOCAL0, "MTTitleCard")
    if err != nil {
        log.Panicf("Failed to connect to syslog: %+v", err)
    }
    defer sl.Close()

    /* Remove date and time flags since the log was already redirected to syslog */
    log.SetFlags(log.Lshortfile)
    log.SetOutput(sl)

    flag.StringVar(&configFile, "config", "", "Path to a config file")
    flag.Parse()
    err = config.Load(configFile)
    if err != nil {
        log.Panicf("Failed to load the configuratin: %+v", err)
    }

    err = page.StartServer(config.Get().Port)
    if err != nil {
        log.Panicf("Failed to start server: %+v", err)
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
    log.Print("Shutting down...")
}
