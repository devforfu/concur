package main

import (
    "../memo"
    "fmt"
    "log"
    "os"
    "time"
)

func main() {
    m := memo.New(memo.HTTPGetBody)
    for url := range incomingURLs() {
        start := time.Now()
        value, err := m.Get(url)
        if err != nil {
            log.Print(err)
        }
        fmt.Printf("%s, %s, %d bytes\n",
            url, time.Since(start), len(value.([]byte)))
    }
}

func incomingURLs() <-chan string {
    ch := make(chan string)
    go func() {
        for _, url := range os.Args[1:] {
            ch <- url
        }
        close(ch)
    }()
    return ch
}