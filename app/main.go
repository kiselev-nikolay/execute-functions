package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("rseDo", rseDo())
	select {}
}

func rseDo() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].String()
		go exec(read(data))
		return nil
	})
}
