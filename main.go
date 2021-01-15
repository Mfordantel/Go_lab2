package main

import (
        "fmt"
        "net/http"
        "io"
        "os"
        "strings"
        "time"
        "bytes"
)

func info(endOfFile bool, buf *bytes.Buffer) {
    for !endOfFile {
        fmt.Println(buf.Len() / 1024, "Kb")
        time.Sleep(time.Second)
    }
}

func main() {
        var url string
        fmt.Print("Enter file URL: ")
        fmt.Scanf("%s\n",&url)

        fileName := url[1+strings.LastIndex(url, "/"):]
        if fileName == "" {
                fmt.Println("Error: Wrong URL format")
                os.Exit(1)
        }

        status, err := http.Get(url)
        if err != nil {
                fmt.Println("Error: Can't get file")
                os.Exit(1)
        }
        defer status.Body.Close()

        file, err := os.Create(fileName)
        if err != nil {
                fmt.Println("Error: Can't create file")
                os.Exit(1)
        }
        defer file.Close()

        var buf bytes.Buffer
        endOfFile := false
        reader := io.TeeReader(status.Body, &buf)
        fmt.Println("Start downloading. Already received:")
        go info(endOfFile, &buf)
        _, err = io.Copy(file, reader)
        if err != nil {
                fmt.Println("Error: Can't copy to file")
                os.Exit(1)
        }
        endOfFile = true

        fmt.Println("Download!")
}
