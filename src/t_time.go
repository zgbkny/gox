package main

import (
    "log"
    "time"
)

var lastTime *time.Time

const bb byte = 0xFF

func main() {
    t := time.Now()
    lastTime = &t
    log.Println(lastTime)
    d := time.Duration(500000000)
    *lastTime = lastTime.Add(d)
    log.Println(lastTime)
    var b byte = 0xFF
    log.Println(b)
    log.Println(bb)
    log.Println(t)
    log.Println(d)
    
    for i := 0; i < 50000000; i++ {
        *lastTime = time.Now()
    }
    log.Println(lastTime)
    log.Println(lastTime.After(time.Now()))
}

