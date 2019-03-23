package main

import (
    "fmt"
    "testing"
    "time"
)

func TestMemo_Get(t *testing.T) {
    cases := []struct{
        url string
        ok bool
    }{
        {"https://gopl.io", true},
        {"https://google.com", true},
        {"https://gopl.io", true},
        {"http://invalid.path", false},
    }
    m := NewMemo(httpGetBody)
    for _, testCase := range cases {
        start := time.Now()
        result, err := m.Get(testCase.url)
        if err != nil && testCase.ok {
            t.Errorf("url failed: %s", testCase.url)
        }
        fmt.Printf("%s, %s, %d",
            testCase.url, time.Since(start), len(result.([]byte)))
    }
}