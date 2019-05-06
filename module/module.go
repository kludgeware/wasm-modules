// +build js

package module

import (
	"fmt"
	"syscall/js"
)

var (
	modules map[string]*jsModule
	debug   bool
)

func init() {
	modules = make(map[string]*jsModule)
}

type Module interface {
	Export(name string, val interface{})
}

type jsModule struct {
	name    string
	exports []js.Value
	debug   bool
}

func (m *jsModule) Export(name string, val interface{}) {
	if debug {
		fmt.Println(name, val)
	}

	js.Global().Get(m.name).Set(name, js.ValueOf(val))
}

func New(name string) Module {
	if _, ok := modules[name]; ok {
		panic("duplicte module found: " + name)
	}

	mod := &jsModule{
		name: name,
	}
	modules[name] = mod

	js.Global().Set(name, make(map[string]interface{}))
	return mod
}
