package memo

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
    "time"
)

type Memo struct {
    f Func
    cache map[string]result
}

type Func func(key string) (interface{}, error)

type result struct {
    value interface{}
    err error
}

func New(f Func) *Memo {
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

func HTTPGetBody(url string) (interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

// Fetch pulls URL strings from urls channel and fetches their content.
func Fetch(w io.Writer, urls <-chan string) {
    m := New(HTTPGetBody)
    var n sync.WaitGroup
    for url := range urls {
        n.Add(1)
        go func(url string) {
            start := time.Now()
            value, err := m.Get(url)
            if err != nil {
                log.Print(err)
            }
            status := fmt.Sprintf("%s, %s, %d bytes",
                url, time.Since(start), len(value.([]byte)))
            _, _ = w.Write([]byte(status))
            n.Done()
        }(url)
    }
    n.Wait()
}