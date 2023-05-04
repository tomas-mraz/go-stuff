package main

import f "github.com/pkg/browser"

// https://pkg.go.dev/github.com/pkg/browser#section-readme

func main() {
	const url = "http://www.seznam.cz/"
	err := f.OpenURL(url)
	if err != nil {
		println("error happen")
	}
}
