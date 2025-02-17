package main

import (
	"riwo/apps"
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func LaunchDefault(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Expected one integer (window id)" // No or too many args
	}
	jsNum := args[0] // Get the js.Value argument

	if jsNum.Type() != js.TypeNumber { // Check if it's a number
		return "Argument must be a number"
	}
	num := jsNum.Int() // Convert js.Value to Go int

	fetchedWindow, ok := wm.AllWindows[strconv.Itoa(num)]
	if !ok {
		// Im really not okay (trust me)
		if wm.Verbose {
			wm.Print("Couldn't start APP_default on window " + strconv.Itoa(num))
		}
		return nil
	}
	apps.APP_default(fetchedWindow)
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	// Print an introductory message to the browser console.
	wm.Print(`
Great, You've found yourself in the console
Then you are likely to want to know this:
- Click LMB to cancel any action
- Press RMB to open context menu
- Select option by pressing RMB
- "New" will open another window after you
  make a selection with RMB
- Select state wants RMB click ("Delete", "Resize")
  or hold ("Move") on desired window
For logging there are:
+ Nothing yet
`)

	wm.Verbose = true // Temporary permanent logging

	wm.AllWindows = make(map[string]*wm.Window)
	wm.ContextMenuHides = make([]js.Value, 0)

	// Set default app for window
	js.Global().Set("LaunchDefault", js.FuncOf(LaunchDefault))
	// Essential for context menu's "New"

	// Window manager core
	wm.InitializeContextMenu()
	wm.InitializeGlobalMouseEvents()

	<-c
}
