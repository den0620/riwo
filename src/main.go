package main

import (
  "fmt"
  "syscall/js"
)

var (
  isDragging   bool
  isMovingMode bool
  startX, startY float64
  activeWindow   js.Value
)

func main() {
  c := make(chan struct{}, 0)

  fmt.Print(`
Great, You've found yourself in the console
Then you are likely to want to know this:
- Press RMB on bg to open menu
- Select option by pressing RMB
- Click LMB to cancel
- Choose window with RMB
- Hold RMB to drag aroung in "move" mode
- Make selection with RMB in "resize" mode
Logging is included
`)

  // Set up the main function for window creation
  js.Global().Set("createDraggableWindow", js.FuncOf(createDraggableWindow))

  // Initialize the context menu
  initializeContextMenu()

  <-c
}

func createDraggableWindow(this js.Value, args []js.Value) interface{} {
  document := js.Global().Get("document")
  body := document.Get("body")

  // Create the window element
  window := document.Call("createElement", "div")
  window.Set("style", "position: absolute; width: 200px; height: 150px; background-color: #f0f0f0; border: solid #55AAAA; padding: 10px;")
  window.Set("innerHTML", "<h3>Draggable Window</h3><p>html p</p>")

  // Set initial position
  window.Get("style").Set("left", "100px")
  window.Get("style").Set("top", "100px")

  // Add window to body
  body.Call("appendChild", window)

  // Prevent default context menu on right-click for the draggable window
  window.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    args[0].Call("stopPropagation")
    fmt.Println("Context menu prevented on draggable window.")
    return nil
  }))

  // Mouse down event for window selection (only if "Move" is active)
  window.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isMovingMode && args[0].Get("button").Int() == 2 { // Only right mouse button (button 2)
      args[0].Call("preventDefault") // Block default right-click behavior
      args[0].Call("stopPropagation")
      activeWindow = window
      startX = args[0].Get("clientX").Float() - window.Get("offsetLeft").Float()
      startY = args[0].Get("clientY").Float() - window.Get("offsetTop").Float()
      isDragging = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-drag.svg), auto") // Drag cursor
      fmt.Println("Dragging initiated.")
    }
    return nil
  }))

  // Global mouse move event
  js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isDragging && isMovingMode && activeWindow.Truthy() {
      x := args[0].Get("clientX").Float() - startX
      y := args[0].Get("clientY").Float() - startY
      activeWindow.Get("style").Set("left", fmt.Sprintf("%fpx", x))
      activeWindow.Get("style").Set("top", fmt.Sprintf("%fpx", y))
      fmt.Println("Window is moving.")
    }
    return nil
  }))

  // Global mouse up event to finalize position
  js.Global().Get("document").Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isDragging {
      isDragging = false
      isMovingMode = false
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto") // Revert cursor to default
      fmt.Println("Dragging ended.")
    }
    return nil
  }))

  return nil
}


func initializeContextMenu() {
  document := js.Global().Get("document")
  body := document.Get("body")

  // Create the context menu
  menu := document.Call("createElement", "div")
  menu.Set("id", "contextMenu")
  menu.Set("style", "position: absolute; display: none; background-color: #EEFFEE; border: solid #8BCE8B; padding: 0;")
  body.Call("appendChild", menu)

  // Add 'Move' option to the context menu
  moveOption := document.Call("createElement", "div")
  moveOption.Set("innerText", "Move")
  moveOption.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto") // Set cursor on hover over "Move" option

  // Set background for hover or selected state
  moveOption.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    moveOption.Get("style").Set("background-color", "#418941")
    moveOption.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
    return nil
  }))

  moveOption.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    moveOption.Get("style").Set("background-color", "#EEFFEE") // Revert to original color
    return nil
  }))

  // Prevent default context menu on right-click for the move option
  moveOption.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    args[0].Call("stopPropagation")
    fmt.Println("Context menu prevented on Move option.")
    return nil
  }))

  // Select 'Move' option
  moveOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 { // Only right mouse button
      args[0].Call("preventDefault") // Block default right-click behavior
      args[0].Call("stopPropagation")
      isMovingMode = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg), auto") // Select cursor
      menu.Get("style").Set("display", "none") // Hide menu after selecting 'Move'
      fmt.Println("Move mode activated.")
    }
    return nil
  }))
  menu.Call("appendChild", moveOption)

  // Hide menu on clicking outside and reset moving mode
  document.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    menu.Get("style").Set("display", "none")
    isMovingMode = false
    js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto") // Revert cursor to default
    fmt.Println("Context menu hidden, Move mode deactivated.")
    return nil
  }))

  // Right-click on background to open the menu (only if not in moving mode)
  body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault") // Prevent the default context menu on right-click
    args[0].Call("stopPropagation")
    if !isMovingMode {
      // Position the custom context menu
      menu.Get("style").Set("left", fmt.Sprintf("%fpx", args[0].Get("clientX").Float()))
      menu.Get("style").Set("top", fmt.Sprintf("%fpx", args[0].Get("clientY").Float()))
      menu.Get("style").Set("display", "block")
      fmt.Println("Context menu displayed.")
    }
    return nil
  }))
}


