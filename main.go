package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type counter struct {
	Seconds int
	Total uint64
}

func (writecount *counter) Write(p []byte) (int, error) {
	n := len(p)
	writecount.Total += uint64(n)
	return n, nil
}

func sec(writecount *counter){
	for {
		time.Sleep(time.Second)
		writecount.Seconds++
		fmt.Println("Time:",writecount.Seconds,"sec")
	}
}

func download(way string, url string) error {

	if !strings.Contains(way, ".") {
		way += ".html"
	}


	out, err := os.Create(way)
	if err != nil {
		return err
	}
	defer out.Close()


	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()


	var buf bytes.Buffer
	endFile := false
	read := io.TeeReader(resp.Body, &buf)
	fmt.Println("Begin download. Already installed:")
	
	count := &counter{}
	go sec(count)

	go func() {
		for !endFile {
			fmt.Println("Downloaded:", buf.Len()/1024, "Kb")
			time.Sleep(time.Second)
		}
	}()

	_, err = io.Copy(out, read)
	endFile = true

	fmt.Println(buf.Len()/1024, "Kb")
	fmt.Println("Downloaded as " + way)
	return err

}
func main() {

	var url string
	fmt.Print("Enter file URL: ")
	fmt.Scanf("%s\n", &url)

	SepUrl := strings.Split(url, "/")

	var file_name string
	for i := 1; file_name == "" && i < len(SepUrl); i++ {
		file_name = SepUrl[len(SepUrl)-i]
	}
        if file_name == "" {
                fmt.Println("Error: Wrong URL format")
                os.Exit(1)
        }

	err := download(file_name, url)
	if err != nil {
                fmt.Println("Error: Can't copy to file")
                os.Exit(1)
	}
}