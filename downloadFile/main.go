package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"time"
)

type progressReader struct {
	reader   io.Reader
	size     int64
	position int64
	start    int64
	context  context.Context
}

func main() {
	tempFile, _ := os.OpenFile("tmp.zip", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer tempFile.Close()

	ctx, cancel := context.WithCancel(context.Background())
	aaa := &progressReader{context: ctx}
	go downloadFile(aaa, ctx, "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-live-server-amd64.iso", tempFile)

	time.Sleep(5 * time.Second)
	aaa.Output()
	time.Sleep(5 * time.Second)
	aaa.Output()

	cancel()
	ctx2, cancel2 := context.WithCancel(context.Background())
	bbb := &progressReader{context: ctx2}
	go downloadFile(bbb, ctx2, "https://releases.ubuntu.com/23.10.1/ubuntu-23.10.1-desktop-amd64.iso", tempFile)
	defer cancel2()

	time.Sleep(5 * time.Second)
	bbb.Output()
	time.Sleep(5 * time.Second)
	bbb.Output()
	fmt.Println("end")
}

func (pr *progressReader) Output() {
	fmt.Print(pr.Progress() + " ")
	fmt.Print(pr.Downloaded() + " ")
	fmt.Print(pr.Speed() + " ")
	fmt.Println(pr.ETA())
}

func (pr *progressReader) ETA() string {
	eta := (pr.size - pr.position) / pr.SpeedNumber()
	finish := time.Now().Add(time.Duration(eta * int64(time.Second)))
	return humanize.RelTime(time.Now(), finish, "", "")
}

// SpeedNumber bytes per second
func (pr *progressReader) SpeedNumber() int64 {
	deltaTime := time.Now().UnixMilli() - pr.start
	return pr.position / (deltaTime / 1000)
}

func (pr *progressReader) Speed() string {
	return humanize.Bytes(uint64(pr.SpeedNumber())) + "/s"
}

func (pr *progressReader) Downloaded() string {
	return humanize.Bytes(uint64(pr.position))
}

func (pr *progressReader) Progress() string {
	percentage := float64(pr.position) / (float64(pr.size) / 100)
	if percentage < 10 {
		return fmt.Sprintf("%2.1f%%", percentage)
	}
	return fmt.Sprintf("%3.0f%%", percentage)
}

func downloadFile(pr *progressReader, ctx context.Context, url string, tempFile *os.File) {
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error during closing connection")
		}
	}(resp.Body)

	pr.reader = resp.Body
	pr.size = resp.ContentLength
	pr.position = 0
	pr.start = time.Now().UnixMilli()

	_, err := io.Copy(tempFile, pr)
	if err != nil {
		fmt.Println("error or cancellation during downloading")
	}
}

func (pr *progressReader) Read(p []byte) (int, error) {
	select {
	case <-pr.context.Done():
		return 0, pr.context.Err()
	default:
		n, err := pr.reader.Read(p)
		if err == nil {
			pr.position += int64(n)
		}
		return n, err
	}
}
