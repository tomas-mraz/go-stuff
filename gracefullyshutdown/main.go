package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("App v2")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Break the Hello loop")
				return
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in a loop")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Break the Hi loop")
				return
			case <-time.After(2 * time.Second):
				fmt.Println("Hi in a loop")
			}
		}
	}()

	wg.Wait()

	fmt.Println("prejmenovat sebe")
	curAppFile := DataFile("aaa.exe")
	newAppFile := DataFile("aaa.exe_new")
	oldAppFile := DataFile("aaa.exe_old")

	fileExist := false
	if _, err := os.Stat(newAppFile); err == nil {
		fileExist = true
	}
	if fileExist {
		os.Rename(curAppFile, oldAppFile)
		os.Rename(newAppFile, curAppFile)
	}

	fmt.Println("Main done")
}
