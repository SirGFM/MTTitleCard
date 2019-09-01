package main

import (
    "log"
)

func setupLogger() {
    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)
}

func cleanLogger() {
}
