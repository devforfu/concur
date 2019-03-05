package thumbnail

import (
    "fmt"
    "path"
    "path/filepath"
    "strings"
)

func DiscoverImages(dirname, formats string) ([]string, error) {
    files := make([]string, 0)
    for _, ext := range SplitPattern(formats) {
        glob := path.Join(dirname, fmt.Sprintf("*.%s", ext))
        matched, err := filepath.Glob(glob)
        if err != nil { return nil, err }
        files = append(files, matched...)
    }
    return files, nil
}

// SplitPattern converts pipe-separated pattern into array of strings.
// For example, the string "jpeg|png" is converted into array ["jpeg", "png"].
func SplitPattern(pattern string) []string {
    result := make([]string, 0)
    for _, p := range strings.Split(pattern, "|") {
        for _, ext := range []string{strings.ToUpper(p), strings.ToLower(p)} {
            result = append(result, ext)
        }
    }
    return result
}