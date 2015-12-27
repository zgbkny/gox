package main

import (
    "log"
    "time"
)

var lastTime *time.Time

func main() {
    t := time.Now()
    lastTime = &t
    d := time.Duration(1000000000)
    log.Println(t)
    log.Println(d)
}

