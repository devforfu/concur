package main

func main() {
    worklist := make(chan []string)
    items := []string{"first", "second", "third"}

    go func() {
        worklist <- items
        close(worklist)
    }()

    seen := make(map[string]bool)
    for list := range worklist {
        for _, item := range list {
            println(item)
            if !seen[item] {
                seen[item] = true
            }
        }
    }
}