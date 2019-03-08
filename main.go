package main

import (
    "./thumbnail"
    "flag"
    "github.com/mitchellh/go-homedir"
    "log"
    "os"
    "path"
    "sync"
    "time"
)

func main() {
    ch := make(chan string)
    go discoverImages(parseFolder(), ch)
    info := makeThumbnails(ch)
    log.Printf("Expected thumbnails count: %d\n", info.Expected)
    log.Printf("Created thumbnails: %d (size: %d)\n", len(info.Files), info.TotalSize)
    log.Printf("Failed: %d\n", info.Failed)
    log.Println("Created files:")
    for i, item := range info.Files {
        log.Printf("\t%d: %s", i, item)
    }
    if info.Failed > 0 {
        for name, err := range info.Errors {
            log.Printf("\tname: %s, error: %s", name, err)
        }
    }
}

type ThumbInfo struct {
    Expected, Failed int
    TotalSize int64
    Files []string
    Errors map[string]string
}

func discoverImages(dirname string, discovered chan<- string) {
    filenames, err := thumbnail.DiscoverImages(dirname, thumbnail.ImageFormats)
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range filenames {
        discovered <- f
    }
}

func makeThumbnails(filenames <-chan string) ThumbInfo {
    type result struct {
        thumbfile string
        size int64
        err error
    }

    var wg sync.WaitGroup
    items := make(chan result)
    maker := thumbnail.NewMaker(thumbnail.JPEG, 128, 128)

    startWait := 1 * time.Millisecond
    wait := startWait

    loop: for {
        select {
        case f := <-filenames:
            wait = startWait
            wg.Add(1)
            log.Printf("Processing file: %s", f)
            go func(f string) {
                defer wg.Done()
                var it result
                if it.thumbfile, it.err = maker.Thumbnail(f); it.err != nil {
                    items <- it
                } else {
                    info, _ := os.Stat(it.thumbfile)
                    it.size = info.Size()
                    items <- it
                }
            }(f)
        default:
            if wait > time.Second {
                log.Printf("Timeout")
                break loop
            } else {
                log.Printf("No files to process. Waiting for %v...", wait)
                time.Sleep(wait)
                wait *= 2
            }
        }
    }

    go func() {
        wg.Wait()
        close(items)
    }()

    var size int64
    total, failed := 0, 0
    thumbs := make([]string, 0)
    errors := make(map[string]string)

    for it := range items {
        total += 1
        if it.err != nil {
            failed += 1
            errors[it.thumbfile] = it.err.Error()
        } else {
            thumbs = append(thumbs, it.thumbfile)
            size += it.size
        }
    }

    return ThumbInfo{Expected:total, Failed:failed, TotalSize:size, Files:thumbs, Errors:errors}
}

func timer(action func ()) time.Duration {
    start := time.Now()
    action()
    elapsed := time.Since(start)
    return elapsed
}

func parseFolder() string {
    defaultFolder, _ := homedir.Dir()
    defaultFolder = path.Join(defaultFolder, "Unsplash")
    folder := flag.String("-d", defaultFolder, "path to the folder with images")
    flag.Parse()
    return *folder
}
