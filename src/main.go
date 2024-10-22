package main

import (
  "fmt"
  "syscall/js"
  "strconv"
)

var (
  isDragging   bool
  isMovingMode bool
  justSelected bool
  startX, startY float64
  activeWindow   js.Value
  ghostWindow    js.Value
  windowCount    int    // Counter for creating multiple windows with unique z-index
  highestZIndex  int = 10 // Track the highest z-index for bringing windows to front
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
- Hold RMB to drag around in "move" mode
- Make selection with RMB in "resize" mode
Logging is included
`)

  js.Global().Set("createDraggableWindow", js.FuncOf(createDraggableWindow))
  initializeContextMenu()

  // Add global mousemove and mouseup listeners only once
  initializeGlobalMouseEvents()

  <-c
}

func initializeGlobalMouseEvents() {
  // Global mouse move event
  js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isDragging && isMovingMode && ghostWindow.Truthy() {
      x := args[0].Get("clientX").Float() - startX
      y := args[0].Get("clientY").Float() - startY
      ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", x))
      ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", y))
      fmt.Println("Ghost window is moving.")
    }
    return nil
  }))

  // Global mouse up event to stop dragging
  js.Global().Get("document").Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isDragging {
      isDragging = false
      isMovingMode = false
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

      // Move the window to the ghost's position
      if ghostWindow.Truthy() && activeWindow.Truthy() {
        activeWindow.Get("style").Set("left", ghostWindow.Get("style").Get("left"))
        activeWindow.Get("style").Set("top", ghostWindow.Get("style").Get("top"))
        ghostWindow.Call("remove") // Remove ghost window
      }

      fmt.Println("Dragging ended and window teleported to ghost position.")
    }
    return nil
  }))
}

func createDraggableWindow(this js.Value, args []js.Value) interface{} {
  document := js.Global().Get("document")
  body := document.Get("body")

  windowCount++ // Increment window counter for unique z-index

  // Create the window element
  window := document.Call("createElement", "div")
  window.Set("style", fmt.Sprintf("position: absolute; z-index: %d; width: 200px; height: 150px; background-color: #f0f0f0; border: solid #55AAAA; padding: 10px;", highestZIndex))
  window.Set("innerHTML", fmt.Sprintf("<h3>Draggable Window %d</h3><p>html p</p>", windowCount))

  // Set initial position
  window.Get("style").Set("left", fmt.Sprintf("%dpx", 100+windowCount*20))
  window.Get("style").Set("top", fmt.Sprintf("%dpx", 100+windowCount*20))

  body.Call("appendChild", window)

  // Prevent the context menu from opening on right-click for windows
  window.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    args[0].Call("stopPropagation")
    fmt.Println("Context menu prevented on draggable window.")
    // Bring window to the front
    highestZIndex++
    window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
    return nil
  }))

  // Mouse down event for selecting and dragging the window (Left-click brings it to front)
  window.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // Left-click to bring window to front
    highestZIndex++
    window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
    fmt.Println("Window brought to front.")

    if isMovingMode && args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      activeWindow = window
      startX = args[0].Get("clientX").Float() - window.Get("offsetLeft").Float()
      startY = args[0].Get("clientY").Float() - window.Get("offsetTop").Float()
      isDragging = true

      // Create ghost window
      ghostWindow = document.Call("createElement", "div")
      rect := window.Call("getBoundingClientRect")
      width := rect.Get("width").Float()
      height := rect.Get("height").Float()

      // Ensure ghost window is above everything during drag
      ghostWindow.Set("style", fmt.Sprintf("position: absolute; z-index: %d; width: %fpx; height: %fpx; border: solid 2px #FF0000; cursor: url(assets/cursor-drag.svg), auto;",
        highestZIndex+1, width, height))
      ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", window.Get("offsetLeft").Float()))
      ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", window.Get("offsetTop").Float()))
      body.Call("appendChild", ghostWindow)

      fmt.Println("Dragging initiated with ghost window.")
    }
    return nil
  }))

  return nil
}


func initializeContextMenu() {
  document := js.Global().Get("document")
  body := document.Get("body")

  // Create the context menu with higher z-index
  menu := document.Call("createElement", "div")
  menu.Set("id", "contextMenu")
  menu.Set("style", "position: absolute; z-index: 999; display: none; background-color: #EEFFEE; border: solid #8BCE8B; padding: 0; text-align: center;")
  body.Call("appendChild", menu)

  // Move option
  moveOption := document.Call("createElement", "div")
  moveOption.Set("innerText", "Move")
  moveOption.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
  moveOption.Get("style").Set("padding", "10px")

  // New option
  newOption := document.Call("createElement", "div")
  newOption.Set("innerText", "New")
  newOption.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
  newOption.Get("style").Set("padding", "10px")

  // Hover and selection effects for Move
  moveOption.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    moveOption.Get("style").Set("background-color", "#418941")
    return nil
  }))
  moveOption.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    moveOption.Get("style").Set("background-color", "#EEFFEE")
    return nil
  }))

  // Hover and selection effects for New
  newOption.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    newOption.Get("style").Set("background-color", "#418941")
    return nil
  }))
  newOption.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    newOption.Get("style").Set("background-color", "#EEFFEE")
    return nil
  }))

  // Move mode activation
  moveOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      isMovingMode = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg), auto")
      menu.Get("style").Set("display", "none")
      fmt.Println("Move mode activated.")
    }
    return nil
  }))

  // New window activation
  newOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      menu.Get("style").Set("display", "none")
      js.Global().Call("createDraggableWindow") // Trigger the creation of a new window
      fmt.Println("New window created.")
    }
    return nil
  }))

  // Add options to the menu
  menu.Call("appendChild", moveOption)
  menu.Call("appendChild", newOption)

  // Global context menu activation
  body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    if !justSelected && !isMovingMode {
      menu.Get("style").Set("left", fmt.Sprintf("%dpx", args[0].Get("clientX").Int()))
      menu.Get("style").Set("top", fmt.Sprintf("%dpx", args[0].Get("clientY").Int()))
      menu.Get("style").Set("display", "block")
    }
    justSelected = false
    return nil
  }))

  // Cancel menu on left-click
  body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 0 {
      menu.Get("style").Set("display", "none")
    }
    return nil
  }))
}

