package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    var files []string

    root := "/root/go/src/github.com/glennswest/libignition"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil {
        panic(err)
    }
    for _, file := range files {
        fmt.Println(file)
    }
}

