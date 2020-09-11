// +build js
// +build wasm

package goo

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/google/uuid"
)

type Canvas struct {
	root *Composite
}

func (c *Canvas) Mount(root *Composite) {
	document = js.Global().Get("document")
	if !document.Truthy() {
		fmt.Println("goo error: Unable to get document object")
		return
	}
	seq = 0
	gooid = uuid.New()
	c.root = root
	rootElem := document.Call("createElement", "div")
	if !rootElem.Truthy() {
		fmt.Println("goo error: Failed to create root element")
		return
	}
	c.root.mountID = "_" + gooid.String() + "_" + strconv.Itoa(seq)
	seq++
	rootElem.Set("id", c.root.mountID)
	rootElem.Set("innerHTML", "")
	document.Get("body").Call("appendChild", rootElem)
	c.root.mount()
	go func() {
		c.root.paint(false)
		dur := 100 * time.Microsecond
		timer := time.NewTimer(dur)
		for {
			<-timer.C
			c.root.reconcile(false)
			timer.Reset(dur)
		}
	}()
	<-make(chan bool)
}

var document js.Value
var seq int
var gooid uuid.UUID

func documentGetElementById(id string) js.Value {
	if !document.Truthy() {
		fmt.Println("goo error: Unable to get document object")
		return js.Null()
	}
	return document.Call("getElementById", id)
}
