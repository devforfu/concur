package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    memo := NewMemo(httpGetBody)
    for url := range incomingURLs() {
        start := time.Now()
        value, err := memo.Get(url)
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

func httpGetBody(url string) (interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

type Memo struct {
    f Func
    cache map[string]result
}

type Func func(key string) (interface{}, error)

type result struct {
    value interface{}
    err error
}

func NewMemo(f Func) *Memo {
    return &Memo{f: f, cache: make(map[string]result)}
}

func (memo *Memo) Get(key string) (interface{}, error) {
    res, ok := memo.cache[key]
    if !ok {
        res.value, res.err = memo.f(key)
        memo.cache[key] = res
    }
    return res.value, res.err
}