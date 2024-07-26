package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func main() {
	fmt.Println("ahoj")
	t := Config{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
