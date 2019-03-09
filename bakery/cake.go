package bakery

import (
    "log"
    "time"
)

type Cake struct { state string }

func (c *Cake) String() string { return c.state }

func Baker(delay time.Duration, cooked chan<- *Cake) {
    for {
        cake := new(Cake)
        log.Println("Baking new cake...")
        time.Sleep(delay)
        cake.state = "cooked"
        cooked <- cake
    }
}

func Icer(delay time.Duration, iced chan<- *Cake, cooked <-chan *Cake) {
    for cake := range cooked {
        log.Println("Icing new cake...")
        time.Sleep(delay)
        cake.state = "iced"
        iced <- cake
    }
}