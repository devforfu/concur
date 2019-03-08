package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    abort := make(chan struct{})
    go func() {
       _, _ = os.Stdin.Read(make([]byte, 1))
       abort <- struct{}{}
    }()
    fmt.Printf("Commencing countdown. Press return to abort.\n")
    select {
    case <-time.After(10*time.Second):
        // Do nothing.
    case <-abort:
        fmt.Println("Launch aborted!")
        return
    }
    launch()
}

func launch() {
    fmt.Println("Launch 🚀")
}