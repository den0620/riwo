package main

import (
	"riwo/wm"
	"syscall/js"
)


func main() {
	wm.RunGuestApp(func(w *wm.RiwoWindow) {
		w.Title = "Monaco Editor (minimal)"
		w.Content.
			Inner("").
			Style("width", "100%").
			Style("height", "100%").
			Style("overflow", "hidden")
		done := make(chan struct{})
		cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			close(done)
			return nil
		})
		defer cb.Release()
		js.Global().Get("__riwoKernel").Call("monacoMountDone", w.Content.DOM(), cb)
		<-done
	})
}
