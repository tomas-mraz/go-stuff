package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
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

	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open("aaa.yaml")
	if err != nil {
		slog.Error("adsdas " + err.Error())
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err2 := d.Decode(&config); err2 != nil {
		slog.Error("asdasd " + err2.Error())
	}

	fmt.Println(config.A)

	/*
		fmt.Println("ahoj")
		t := Config{}

		err := yaml.Unmarshal([]byte(data), &t)
		if err != nil {
			slog.Error("error: ")
		}
		fmt.Printf("--- t:\n%v\n\n", t)

		d :=
		d, err2 = yaml.Marshal(&t)
		if err2 != nil {
			slog.Error("error2: ")
		}
	*/

	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
