package memo_test

import (
    "../memo"
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
    m := memo.New(memo.HTTPGetBody)
    for _, testCase := range cases {
        start := time.Now()
        value, err := m.Get(testCase.url)
        if err != nil {
            if testCase.ok {
                t.Errorf("url failed: %s\n", testCase.url)
            }
        } else {
            fmt.Printf("%s, %s, %d\n",
                testCase.url, time.Since(start), len(value.([]byte)))
        }
    }
}