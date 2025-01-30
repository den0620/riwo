package main


import (
	"syscall/js"
	"riwo/wm"
  "riwo/fs"
  "riwo/rc"
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
For logging there are:
"wm_logging()"
"fs_logging()" (not yet)
"rc_logging()" (not yet)
`)

  js.Global().Set("wm_logging", js.FuncOf(wm.GoVerbose))
  js.Global().Set("fs_logging", js.FuncOf(fs.GoVerbose))
  js.Global().Set("rc_logging", js.FuncOf(rc.GoVerbose))

  // Not implemented yet
  //fs.InitializeStructure()

  // Generate menu
  wm.InitializeContextMenu()
  // Add global mousemove and mouseup listeners only once
  wm.InitializeGlobalMouseEvents()

  // Not implemented yet
  //rc.StartInWindow(1)

  <-c
}

