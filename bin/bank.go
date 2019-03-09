package main

import (
    "../bank"
    "fmt"
)

func main() {
    fmt.Printf("Balance: %d\n", bank.Balance())
    bank.Deposit(9000)
    fmt.Printf("After update: %d\n", bank.Balance())
}