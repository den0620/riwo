package wm

import "syscall/js"

// RunGuestApp runs constructor with window content bound to __riwoGuestBootstrap.pane from the JS kernel
// Blocks forever like the wm main goroutine pattern.
func RunGuestApp(run func(*RiwoWindow)) {
	c := make(chan struct{})
	b := js.Global().Get("__riwoGuestBootstrap")
	pane := b.Get("pane")
	wid := b.Get("windowId").Int()
	pj := pane
	content := CreateFrom(&pj)
	rw := &RiwoWindow{
		ID:      wid,
		Content: content,
		Title:   "",
	}
	run(rw)
	<-c
}
