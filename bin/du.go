package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
    "time"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
    flag.Parse()
    roots := flag.Args()
    if len(roots) == 0 {
        roots = []string{"."}
    }

    fileSizes := make(chan int64)
    var n sync.WaitGroup
    for _, root := range roots {
        n.Add(1)
        go walkDir(root, &n, fileSizes)
    }

    go func() {
        n.Wait()
        close(fileSizes)
    }()

    var tick <-chan time.Time
    if *verbose {
        tick = time.Tick(500 * time.Millisecond)
    }

    var nFiles, nBytes int64
    loop: for {
        select {
        case size, ok := <-fileSizes:
            if !ok {
                break loop
            }
            nFiles++
            nBytes += size
        case <-tick:
            printDiskUsage(nFiles, nBytes)
        }
    }

    printDiskUsage(nFiles, nBytes)
}

func printDiskUsage(nFiles, nBytes int64) {
    fmt.Printf("%d files  %.1f GB\n", nFiles, float64(nBytes)/1e9)
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
    defer n.Done()
    for _, entry := range dirents(dir) {
        if entry.IsDir() {
            n.Add(1)
            subdir := filepath.Join(dir, entry.Name())
            go walkDir(subdir, n, fileSizes)
        } else {
            fileSizes <- entry.Size()
        }
    }
}

type Semaphore struct {
    Size int
    tokens chan struct{}
}
func NewSemaphore(size int) *Semaphore {
    s := Semaphore{Size:size}
    s.tokens = make(chan struct{}, size)
    return &s
}
func (s *Semaphore) Acquire() { s.tokens <- struct{}{} }
func (s *Semaphore) Release() { <-s.tokens }

var sema = NewSemaphore(20)

func dirents(dir string) []os.FileInfo {
    sema.Acquire()
    defer sema.Release()
    entries, err := ioutil.ReadDir(dir)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "du: %v\n", err)
        return nil
    }
    return entries
}