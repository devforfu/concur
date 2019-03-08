package main

import (
    "./thumbnail"
    "flag"
    "github.com/mitchellh/go-homedir"
    "log"
    "path"
    "time"
)

func main() {
    makeThumbnails(parseFolder())
}

func makeThumbnails(dirname string) {
    log.Printf("Converting files from the folder: %s", dirname)

    maker := thumbnail.NewMaker(thumbnail.JPEG, 128, 128)
    filenames, err := thumbnail.DiscoverImages(dirname, thumbnail.ImageFormats)
    if err != nil { log.Fatal(err) }

    elapsed := timer(func() {
        ch := make(chan struct{})
        for _, f := range filenames {
            go func(f string) {
                if outfile, err := maker.Thumbnail(f); err != nil {
                    log.Println(err)
                } else {
                    log.Printf("Created file: %s", outfile)
                }
                ch <- struct{}{}
            }(f)
        }
        for range filenames {
            <-ch
        }
    })

    log.Printf("Elapsed time: %v", elapsed)
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
