package main

import (
    "log"
    "log/syslog"
)

var sl *syslog.Writer

func setupLogger() {
    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)

    var err error
    sl, err = syslog.New(syslog.LOG_LOCAL0, "MTTitleCard")
    if err != nil {
        log.Panicf("Failed to connect to syslog: %+v", err)
    }

    /* Remove date and time flags since the log was already redirected to syslog */
    log.SetFlags(log.Lshortfile)
    log.SetOutput(sl)
}

func cleanLogger() {
    sl.Close()
}
