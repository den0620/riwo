package wm

import "syscall/js"

// RunGuestApp consumes the one pending Go-guest bootstrap latch from the JS kernel, then runs the app.
func RunGuestApp(run func(*RiwoWindow)) {
	k := js.Global().Get("__riwoKernel")
	if !k.Truthy() {
		return
	}
	boot := k.Call("consumeGuestBootstrap")
	if !boot.Truthy() || boot.Type() == js.TypeNull || boot.Type() == js.TypeUndefined {
		return
	}
	wid := boot.Get("windowId").Int()
	pane := boot.Get("pane")
	pj := pane
	content := CreateFrom(&pj)
	rw := &RiwoWindow{
		ID:      wid,
		Content: content,
		Title:   "",
	}
	run(rw)
	select {}
}
