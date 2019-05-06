// +build js

package main

import (
	"fmt"
	"syscall/js"
	"strings"
	"time"
	"github.com/kludgeware/wasm-modules/module"
	"github.com/nilslice/protolock"
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

func tick(this js.Value, args []js.Value) interface{} {
	go func() {
		t := time.NewTicker(time.Second * 2)
		for {
			fmt.Println(<-t.C)
		}
	}()
	return nil
}

func status(this js.Value, args []js.Value) interface{} {
	return nil
}

func document() js.Value {
	return js.Global().Get("document")
}

func querySelector(selector string) js.Value {
	return document().Call("querySelector", selector)
}

func createElement(element string) js.Value {
	return document().Call("createElement", element)
}

func bind(el js.Value, event string, cb js.Func) {
	el.Call("addEventListener", event, cb)
}

func main() {
	c := make(chan struct{})

	mod := module.New("my_module")

	mod.Export("toGo", js.FuncOf(toGo))
	mod.Export("add", js.FuncOf(add))

	mod2 := module.New("another")
	mod2.Export("timeNow", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return time.Now().String()
	}))

	mod3 := module.New("timer")
	mod3.Export("tick", js.FuncOf(tick))

	body := querySelector("body")

	inputLock := createElement("textarea")
	inputLock.Set("placeholder", "add your proto.lock file here...")
	inputProtos := createElement("textarea")
	inputProtos.Set("placeholder", "add your proto here...")
	
	submit := createElement("input")
	submit.Set("value", "Check")
	submit.Set("type", "submit")

	bind(submit, "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		current, err := protolock.FromReader(strings.NewReader(inputLock.Get("value").String()))
		if err != nil {
			fmt.Println("bind submit click:", err)
			return nil 
		}
		fmt.Println("current:", current)

		entry, err := protolock.Parse("input_proto", strings.NewReader(inputProtos.Get("value").String()))
		if err != nil {
			fmt.Println("bind submit click:", err)
			return nil 
		}
		updated := protolock.Protolock{
			Definitions: append([]protolock.Definition{}, 
				protolock.Definition{
					Filepath: current.Definitions[0].Filepath,
					Def: entry,
				},
			),
		}
		fmt.Println("updated:", updated)

		report, err := protolock.Compare(current, updated)
		if err != nil {
			if err == protolock.ErrWarningsFound {
				for _, warning := range report.Warnings {
					fmt.Println(warning)
				}
			}

			fmt.Println("bind submit click:", err)
			return nil 
		}

		return nil
	}))
	
	body.Call("appendChild", submit)
	body.Call("appendChild", inputLock)
	body.Call("appendChild", inputProtos)

	fmt.Println("loaded")

	<-c
}

