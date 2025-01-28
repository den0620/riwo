package main


import (
	"syscall/js"
	"riwo/wm"
)


func main() {
  c := make(chan struct{}, 0)

  wm.Print(`
Great, You've found yourself in the console
Then you are likely to want to know this:
- Press RMB to open context menu
- Select option by pressing RMB
- Click LMB to cancel
- "New" will open another window after you
  make a selection with RMB
- Choose window with RMB
- "Delete" will remove selected window
- Hold RMB to drag around in "Move" mode
- Make selection with RMB in "Resize" mode
For logging type "logging()"
`)

  js.Global().Set("logging", js.FuncOf(wm.GoVerbose))

  // Generate menu
  wm.InitializeContextMenu()

  // Add global mousemove and mouseup listeners only once
  wm.InitializeGlobalMouseEvents()

  <-c
}

