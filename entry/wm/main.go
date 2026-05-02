package main

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
)


func logging(this js.Value, args []js.Value) interface{} {
	wm.Verbose = !wm.Verbose
	if wm.Verbose {
		js.Global().Get("console").Call("log", "Logging is now ON")
	} else {
		js.Global().Get("console").Call("log", "Logging is now OFF")
	}
	return nil
}

func launchDefault(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Expected one integer (window id)"
	}
	jsNum := args[0]
	if jsNum.Type() != js.TypeNumber {
		return "Argument must be a number"
	}
	num := jsNum.Int()
	fetchedWindow, ok := wm.AllWindows[strconv.Itoa(num)]
	if !ok {
		js.Global().Get("console").Call("log", "Couldn't start launchpad on window "+strconv.Itoa(num))
		return nil
	}
	wm.OpenLaunchpad(fetchedWindow)
	return nil
}

func main() {
	c := make(chan struct{})

	js.Global().Get("console").Call("log", `
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
+ Logging()
`)

	js.Global().Set("Logging", js.FuncOf(logging))

	wm.AllWindows = make(map[string]*wm.RiwoWindow)
	wm.ContextMenuHides = make([]js.Value, 0)

	js.Global().Set("LaunchDefault", js.FuncOf(launchDefault))

	wm.InitializeContextMenu()
	wm.InitializeGlobalMouseEvents()

	<-c
}
