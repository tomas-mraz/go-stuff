package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	downloadFile("https://wordpress.org/wordpress-4.4.2.zip", "a.zip")
}

func downloadFile(url string, filePath string) {
	localFile, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer localFile.Close()

	name := path.Base(url)
	fmt.Println("downloading file", name, "from", url)

	// get size
	headResp, _ := http.Head(url)
	defer headResp.Body.Close()
	size, _ := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)
	done := make(chan int64)
	go printDownloadPercent(done, size)

	// download
	start := time.Now()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	n, _ := io.Copy(localFile, resp.Body)
	done <- n

	elapsed := time.Since(start)
	fmt.Printf("Download completed in %s\n", elapsed)
}

func printDownloadPercent(done chan int64, total int64) {
	var stop bool = false
	for {
		select {
		case <-done:
			stop = true
		default:
			progress := <-done
			fmt.Println(progress)
			//fmt.Println(fmt.Sprintf("%.0f", progress/total*100), "%")
		}

		if stop {
			break
		}
		//time.Sleep(10 * time.Millisecond)
	}
}
