// +build js
// +build wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/joshuanario/goo"
)

// First, make sure the javascript glue code is served with "cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./test/assets"
// Then compile with "GOOS=js GOARCH=wasm go build -o ./test/assets/goo_test.wasm ./test/wasm/main.go"
func main() {
	fmt.Println("bootstrapping goo test app")
	document := js.Global().Get("document")
	if !document.Truthy() {
		fmt.Println("goo error: Unable to get document object")
		return
	}
	timers := goo.Composite{
		InitialState: "I live now...at exactly " + time.Now().String(),
		Components:   nil,
		HTML: func(state interface{}) string {
			msg := ""
			for i := 0; i < 20; i++ {
				msg = msg + "<h2>" + state.(string) + "</h2>"
			}
			return msg
		},
	}
	dur := 10 * time.Millisecond
	wait := time.NewTimer(dur)
	go func() {
		for {
			timers.SetState("I live now...at exactly " + time.Now().String())
			<-wait.C
			wait.Reset(dur)
		}
	}()
	clickID := "click-me"
	var onClickCounter js.Func
	clickCounter := goo.Composite{
		InitialState: 0,
		Components:   nil,
		HTML: func(state interface{}) string {
			return fmt.Sprintf("<button id=\"%s\">Click me</button><h1>Click counter: %d</h1>", clickID, state.(int))
		},
		BeforePaint: func() {
			elem := document.Call("getElementById", clickID)
			if elem.Truthy() {
				elem.Call("removeEventListener", "click", onClickCounter)
			}
		},
		AfterPaint: func() {
			elem := document.Call("getElementById", clickID)
			if elem.Truthy() {
				elem.Call("addEventListener", "click", onClickCounter)
			}
		},
	}
	onClickCounter = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		clickCounter.SetState(clickCounter.GetState().(int) + 1)
		return onClickCounter
	})
	root := goo.Composite{
		Components: []*goo.Composite{
			&clickCounter,
			&timers,
		},
	}
	var canvas goo.Canvas
	canvas.Mount(&root)
}
