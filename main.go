package main


import (
  "syscall/js"
  "riwo/wm"
  "riwo/demo"
)


func colorsEntry(this js.Value, args []js.Value) interface{} {
  if len(args) != 1 {
    return "Expected one integer (window id)" // Return a value to JavaScript if needed
  }

  jsNum := args[0] // Get the js.Value argument

  if jsNum.Type() != js.TypeNumber { // Check if it's a number (or string that can be parsed)
    return "Argument must be a number"
  }

  num := jsNum.Int() // Convert js.Value to Go int

  demo.Colors_StartInWindow(num)
  return nil
}


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
`)

  js.Global().Set("wm_logging", js.FuncOf(wm.GoVerbose))

  // colors entry
  js.Global().Set("colors", js.FuncOf(colorsEntry))

  // Generate menu
  wm.InitializeContextMenu()
  // Add global mousemove and mouseup listeners only once
  wm.InitializeGlobalMouseEvents()

  <-c
}

