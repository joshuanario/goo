// +build js
// +build wasm

package goo

import (
	"strconv"
)

type Composite struct {
	mountID      string
	oldStatePtr  *interface{}
	currStatePtr *interface{}
	InitialState interface{}
	Components   []*Composite
	HTML         func(state interface{}) string
	BeforePaint  func()
	AfterPaint   func()
}

func (c *Composite) GetState() interface{} {
	return *(c).currStatePtr
}

func (c *Composite) SetState(newState interface{}) {
	c.currStatePtr = &newState
}

func (parent *Composite) mount() {
	parent.SetState(parent.InitialState)
	parentElem := documentGetElementById(parent.mountID)
	for _, child := range parent.Components {
		chidlElem := document.Call("createElement", "div")
		child.mountID = "_" + gooid.String() + "_" + strconv.Itoa(seq)
		seq++
		child.SetState(child.InitialState)
		chidlElem.Set("id", child.mountID)
		chidlElem.Set("innerHTML", "")
		parentElem.Call("appendChild", chidlElem)
	}
}

func (c *Composite) reconcile(isReconciled bool) {
	willReconcile := false
	hasStateDiff := c.oldStatePtr != c.currStatePtr
	if hasStateDiff && !isReconciled {
		willReconcile = true
		c.paint(false)
	}
	c.oldStatePtr = c.currStatePtr
	if len(c.Components) <= 0 {
		return
	}
	for _, i := range c.Components {
		i.reconcile(willReconcile)
	}
}

func (c *Composite) paint(isDOMUpdated bool) {
	if c.BeforePaint != nil {
		c.BeforePaint()
	}
	if len(c.Components) > 0 {
		for _, i := range c.Components {
			if i.BeforePaint != nil {
				i.BeforePaint()
			}
		}
	}
	mount := documentGetElementById(c.mountID)
	mount.Set("innerHTML", c.innerHTML())
	if len(c.Components) > 0 {
		for _, i := range c.Components {
			if i.AfterPaint != nil {
				i.AfterPaint()
			}
		}
	}
	if c.AfterPaint != nil {
		c.AfterPaint()
	}
}

func (c *Composite) innerHTML() string {
	var ret string
	if len(c.Components) <= 0 {
		return c.HTML(*(c.currStatePtr))
	}
	for _, i := range c.Components {
		chidlElem := document.Call("createElement", "div")
		chidlElem.Set("id", i.mountID)
		chidlElem.Set("innerHTML", i.innerHTML())
		ret += chidlElem.Get("outerHTML").String()
	}
	return ret
}
