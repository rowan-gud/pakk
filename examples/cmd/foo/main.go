package main

import (
	_ "embed"
	"fmt"

	"github.com/rowan-gud/pakk/examples/pkga"
)

//go:embed NAME
var name []byte

func main() {
	greeting := pkga.Greet(string(name))

	fmt.Println(greeting)
}
