// +build js

package main

import (
	"fmt"
	"syscall/js"
	"github.com/kludgeware/wasm-modules/module"
)

func toGo(this js.Value, args []js.Value) interface{} {
	fmt.Println(this, args)
	return nil
}

func add(this js.Value, args []js.Value) interface{} {
	var sum int
	for i := range args {
		sum+=args[i].Int()
	}

	return sum
}

func main() {
	c := make(chan struct{})

	mod := module.New("my_module")

	mod.Export("toGo", js.FuncOf(toGo))
	mod.Export("add", js.FuncOf(add))

	fmt.Println("loaded")
	<-c
}

