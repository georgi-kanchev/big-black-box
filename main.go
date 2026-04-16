package main

import (
	"big-black-box/engine"
	"big-black-box/window"
	"fmt"
)

func main() {
	engine.Run(func() {
		fmt.Printf("%v\n", window.IsFocused())
	})
}
