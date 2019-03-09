package main

import (
    "../bakery"
    "bufio"
    "log"
    "os"
    "time"
)

func main() {
    cooked := make(chan *bakery.Cake)
    iced := make(chan *bakery.Cake)
    stop := make(chan struct{}, 1)

    go bakery.Baker(time.Second, cooked)
    go bakery.Icer(1500 * time.Millisecond, iced, cooked)
    go func(stop chan struct{}) {
        scan := bufio.NewScanner(os.Stdin)
        if scan.Scan() {
            stop <- struct{}{}
        }
    }(stop)

    for {
        select {
        case cake := <- iced:
            log.Printf("\tCake is ready: %s\n", cake.String())
        case <-stop:
            log.Println("Ended!")
            break
        }
    }
}

func init() {
    log.SetFlags(0)
}