/*
Application default reads AppRegistry
and enumerates through all apps that registered themselves
*/

package apps

import (
	"riwo/wm"
	"syscall/js"
)

func init() {
	// Register the default app itself.
	AppRegistry["Default"] = APP_default
}

func APP_default(window *wm.Window) {
	document := js.Global().Get("document")

	// Create a container div for the grid
	container := document.Call("createElement", "div")
	container.Get("style").Set("display", "grid")
	container.Get("style").Set("gridTemplateColumns", "repeat(auto-fit, minmax(100px, 1fr))")
	container.Get("style").Set("gap", "10px")
	container.Get("style").Set("padding", "10px")

	// Iterate over AppRegistry and create a button for each app (skip Default to avoid recursion)
	for appName, appFunc := range AppRegistry {
		// Temporary, in future make new design for APP_default
		btn := wm.CreateMenuOption("APP_" + appName)
		// On click, call the corresponding app's main function.
		btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// When an app button is clicked, transfer the window to that app.
			wm.Print("App " + appName + " selected")
			appFunc(window)
			return nil
		}))
		container.Call("appendChild", btn)
	}

	// Clear existing content and show the grid.
	window.Element.Set("innerHTML", "")
	window.Element.Call("appendChild", container)

	// Register a context menu entry so you can return to the grid from any app.
	// This assumes wm.CreateMenuOption creates a js.Value representing a menu entry.
	defaultMenu := wm.CreateMenuOption("Apps")
	defaultMenu.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		APP_default(window)
		return nil
	}))
	return
}
