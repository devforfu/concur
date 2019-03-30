package memotest

import (
    "../memo"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "testing"
    "time"
)

func incomingURLs() <-chan string {
    ch := make(chan string)
    go func() {
        for _, url := range []string{
            "https://golang.org",
            "https://godoc.org",
            "https://play.golang.org",
            "http://gopl.io",
            "https://golang.org",
            "https://godoc.org",
            "https://play.golang.org",
            "http://gopl.io",
        } {
            ch <- url
        }
        close(ch)
    }()
    return ch
}

func httpGetBody(url string) (interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

var HTTPGetBody = httpGetBody

func Sequential(t *testing.T) {
    m := memo.New(HTTPGetBody)
    for url := range incomingURLs() {
        start := time.Now()
        value, err := m.Get(url)
        if err != nil {
            log.Print(err)
            continue
        }
        fmt.Printf("%s, %s, %d bytes\n",
            url, time.Since(start), len(value.([]byte)))
    }
}