package main

import (
  "fmt"
  "syscall/js"
  "strconv"
)

var (
  isDragging     bool
  isMovingMode   bool
  isResizingMode bool
  isResizingInit bool
  justSelected   bool
  isDeleteMode   bool
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
- Press RMB on background to open context menu
- Select option by pressing RMB
- "new" will open another window
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
  js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // Global mouse move event
    if isDragging && isMovingMode && ghostWindow.Truthy() {
      x := args[0].Get("clientX").Float() - startX
      y := args[0].Get("clientY").Float() - startY
      ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", x))
      ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", y))
      fmt.Println("Ghost window is moving.")
    }

    // Global mouse move event to adjust ghost window size during resizing
    if ghostWindow.Truthy() && isResizingMode && isResizingInit && isDragging && !isMovingMode {
      currentX := args[0].Get("clientX").Float()
      currentY := args[0].Get("clientY").Float()

      // Calculate and update ghost window size and position based on selection
      width := currentX - startX
      height := currentY - startY
      ghostWindow.Get("style").Set("width", fmt.Sprintf("%fpx", width))
      ghostWindow.Get("style").Set("height", fmt.Sprintf("%fpx", height))

      // Handle direction to ensure ghost window moves according to selection direction
      if width < 0 {
        ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", currentX))
        ghostWindow.Get("style").Set("width", fmt.Sprintf("%fpx", -width))
      }
      if height < 0 {
        ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", currentY))
        ghostWindow.Get("style").Set("height", fmt.Sprintf("%fpx", -height))
      }
      fmt.Println("Ghost window is resizing with freeform selection.")
    }
    return nil
  }))

  js.Global().Get("document").Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // Global mouse up event to stop dragging
    if isMovingMode && isDragging {
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

    // Global mouse up event to finalize resizing and apply to window
    if ghostWindow.Truthy() && activeWindow.Truthy() && isResizingMode && isResizingInit && isDragging && !isMovingMode {
      // Apply the ghost window's size and position to the actual window
      activeWindow.Get("style").Set("left", ghostWindow.Get("style").Get("left"))
      activeWindow.Get("style").Set("top", ghostWindow.Get("style").Get("top"))
      activeWindow.Get("style").Set("width", ghostWindow.Get("style").Get("width"))
      activeWindow.Get("style").Set("height", ghostWindow.Get("style").Get("height"))
      ghostWindow.Call("remove") // Remove the ghost window
      ghostWindow = js.Null() // Reset ghost window reference
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
      isResizingMode = false
      isResizingInit = false
      isDragging = false

      fmt.Println("Resizing completed and window resized to match selection.")
    }
    return nil
  }))

  // second rmb
  js.Global().Get("document").Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if isResizingMode && args[0].Get("button").Int() == 2 && isResizingInit {
      args[0].Call("preventDefault")
      // Second RMB hold - Start resizing by creating a selection anywhere
      fmt.Println("Second right-click: Resizing initiated.")
      isDragging = true
      // Reset selection start position
      startX = args[0].Get("clientX").Float()
      startY = args[0].Get("clientY").Float()

      // Create ghost window for resizing
      ghostWindow = js.Global().Get("document").Call("createElement", "div")
      ghostWindow.Set("style", fmt.Sprintf("position: absolute; z-index: %d; border: solid 2px #FF0000;", highestZIndex+1))
      ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", startX))
      ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", startY))
      js.Global().Get("document").Get("body").Call("appendChild", ghostWindow)

      fmt.Println("Resizing initiated with freeform selection.")
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
    //args[0].Call("stopPropagation")
    fmt.Println("Caught click on window")
    // Bring window to the front
    if !isResizingInit {
      highestZIndex++
      window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
    }
    return nil
  }))

  window.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // Right-click (RMB) on the window to select it for resizing, second right-click activates resizing
    if isResizingMode && !isResizingInit && args[0].Get("button").Int() == 2 {
      // First RMB hold - Select the window for resizing
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      fmt.Println("First right-click: Window selected for resizing.")
      highestZIndex++
      window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
      activeWindow = window
      isResizingInit = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-selection.svg), auto")
    }

    // Mouse down event for selecting and dragging the window (Left-click brings it to front)
    if !isResizingInit {
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
        ghostWindow.Set("style", fmt.Sprintf("position: absolute; z-index: %d; width: %fpx; height: %fpx; border: solid 2px #FF0000; cursor: url(assets/cursor-drag.svg), auto;", highestZIndex+1, width, height))
        ghostWindow.Get("style").Set("left", fmt.Sprintf("%fpx", window.Get("offsetLeft").Float()))
        ghostWindow.Get("style").Set("top", fmt.Sprintf("%fpx", window.Get("offsetTop").Float()))
        body.Call("appendChild", ghostWindow)

        fmt.Println("Dragging initiated with ghost window.")
      }
    }

    // Right-click (RMB) deletes the window in delete mode
    if isDeleteMode && args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      window.Call("remove") // Delete the window
      isDeleteMode = false
      justSelected = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
      fmt.Println("Window deleted.")
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
  menu.Set("style", "position: absolute; display: none; background-color: #EEFFEE; border: solid #8BCE8B; padding: 0; text-align: center;")
  body.Call("appendChild", menu)

  // Move, New, Resize, and Delete options
  moveOption := createMenuOption("Move")
  newOption := createMenuOption("New")
  resizeOption := createMenuOption("Resize")
  deleteOption := createMenuOption("Delete")

  menu.Call("appendChild", moveOption)
  menu.Call("appendChild", newOption)
  menu.Call("appendChild", resizeOption)
  menu.Call("appendChild", deleteOption)

  // Cancel menu and actions on left-click
  body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 0 {
      menu.Get("style").Set("display", "none")
      isDeleteMode = false // Cancel delete mode on left-click
      isResizingMode = false
      isResizingInit = false
      isDragging = false
      isMovingMode = false
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
    }
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

  // Resize mode activation (second RMB hold will activate actual resizing)
  resizeOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      isResizingMode = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg), auto")
      menu.Get("style").Set("display", "none")
      fmt.Println("Resize mode activated.")
    }
    return nil
  }))

  // Delete mode activation
  deleteOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      isDeleteMode = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg), auto")
      menu.Get("style").Set("display", "none")
      fmt.Println("Delete mode activated.")
    }
    return nil
  }))

  // Global context menu activation
  body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    if !justSelected && !isMovingMode && !isResizingMode && !isDeleteMode{
      // Adjust z-index dynamically based on highestZIndex
      menu.Get("style").Set("z-index", strconv.Itoa(highestZIndex+10))
      menu.Get("style").Set("left", fmt.Sprintf("%dpx", args[0].Get("clientX").Int()))
      menu.Get("style").Set("top", fmt.Sprintf("%dpx", args[0].Get("clientY").Int()))
      menu.Get("style").Set("display", "block")
    }
    justSelected = false
    return nil
  }))
}

// Create a helper function for menu option creation
func createMenuOption(optionText string) js.Value {
  document := js.Global().Get("document")
  option := document.Call("createElement", "div")
  option.Set("innerText", optionText)
  option.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
  option.Get("style").Set("padding", "10px")

  // Hover and selection effects
  option.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    option.Get("style").Set("background-color", "#418941")
    return nil
  }))
  option.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    option.Get("style").Set("background-color", "#EEFFEE")
    return nil
  }))
  return option
}

