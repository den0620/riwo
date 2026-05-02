package wm

import "syscall/js"

// KernelSpawnGuest clears the pane and asks the JS kernel to load the named guest wasm into this window
func KernelSpawnGuest(appName string, window *RiwoWindow) {
	js.Global().Get("__riwoKernel").Call("spawnGuestApp", appName, window.ID, window.Content.DOM())
}
