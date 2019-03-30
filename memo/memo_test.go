package memo_test

import (
    "../memotest"
    "testing"
)

func TestSequential(t *testing.T) {
    memotest.Sequential(t)
}

func TestConcurrent(t *testing.T) {
    memotest.Concurrent(t)
}

