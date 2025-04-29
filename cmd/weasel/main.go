package main

import (
	"os"

	"github.com/AvoidMe/weasel/pkg/compiler"
)

func main() {
	source, err := os.ReadFile("examples/tmp.wsl")
	if err != nil {
		panic(err)
	}

	program, err := compiler.Compile(source)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("weasel_otput.go", program, 0666)
	if err != nil {
		panic(err)
	}
}
