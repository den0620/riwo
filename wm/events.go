package wm


import (
  "syscall/js"
  "strconv"
)


func InitializeGlobalMouseEvents() {
  js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // Global mouse move event
    if isDragging && isMovingMode && ghostWindow.Truthy() {
      x := args[0].Get("clientX").Float() - startX
      y := args[0].Get("clientY").Float() - startY
      ghostWindow.Get("style").Set("left", Ftoa(x) + "px")
      ghostWindow.Get("style").Set("top", Ftoa(y) + "px")
      if verbose {Print("Ghost window is moving.")}
    }

    // Global mouse move event to adjust ghost window size during resizing
    if ghostWindow.Truthy() && isResizingMode && isResizingInit && isDragging && !isMovingMode {
      currentX := args[0].Get("clientX").Float()
      currentY := args[0].Get("clientY").Float()

      // Calculate and update ghost window size and position based on selection
      width := currentX - startX
      height := currentY - startY
      ghostWindow.Get("style").Set("width", Ftoa(width) + "px")
      ghostWindow.Get("style").Set("height", Ftoa(height) + "px")

      // Handle direction to ensure ghost window moves according to selection direction
      if width < 0 {
        ghostWindow.Get("style").Set("left", Ftoa(currentX) + "px")
        ghostWindow.Get("style").Set("width", Ftoa(-width) + "px")
      }
      if height < 0 {
        ghostWindow.Get("style").Set("top", Ftoa(currentY) + "px")
        ghostWindow.Get("style").Set("height", Ftoa(-height) + "px")
      }
      if verbose {Print("Ghost window is resizing with freeform selection.")}
    }

    if ghostWindow.Truthy() && isNewMode && isDragging {
      currentX := args[0].Get("clientX").Float()
      currentY := args[0].Get("clientY").Float()

      // Calculate and update ghost window size and position based on selection
      width := currentX - startX
      height := currentY - startY
      ghostWindow.Get("style").Set("width", Ftoa(width) + "px")
      ghostWindow.Get("style").Set("height", Ftoa(height) + "px")

      // Handle direction to ensure ghost window moves according to selection direction
      if width < 0 {
        ghostWindow.Get("style").Set("left", Ftoa(currentX) + "px")
        ghostWindow.Get("style").Set("width", Ftoa(-width) + "px")
      }
      if height < 0 {
        ghostWindow.Get("style").Set("top", Ftoa(currentY) + "px")
        ghostWindow.Get("style").Set("height", Ftoa(-height) + "px")
      }
      if verbose {Print("New window selection resizing.")}
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
      justSelected = false

      if verbose {Print("Dragging ended and window teleported to ghost position.")}
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
      justSelected = false
      isDragging = false

      if verbose {Print("Resizing completed and window resized to match selection.")}
    }

    if isNewMode && isDragging {
      isNewMode = false
      isDragging = false
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

      // Create a new window at the ghost window's position and size
      if ghostWindow.Truthy() {
        x := ghostWindow.Get("style").Get("left").String()
        y := ghostWindow.Get("style").Get("top").String()
        width := ghostWindow.Get("style").Get("width").String()
        height := ghostWindow.Get("style").Get("height").String()
        ghostWindow.Call("remove")
        ghostWindow = js.Null()

        CreateDraggableWindow(x, y, width, height)

        if verbose {Print("New window created at selected area.")}
      }
    }
    return nil
  }))

  // second rmb
  js.Global().Get("document").Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if (isResizingMode || isNewMode) && args[0].Get("button").Int() == 2 && (isResizingInit || isNewMode) {
      args[0].Call("preventDefault")
      // Second RMB hold - Start resizing by creating a selection anywhere
      if verbose {Print("Second right-click: Resizing initiated.")}
      isDragging = true
      // Reset selection start position
      startX = args[0].Get("clientX").Float()
      startY = args[0].Get("clientY").Float()

      // Create ghost window for resizing
      ghostWindow = js.Global().Get("document").Call("createElement", "div")
      ghostWindow.Set("style", "position: absolute; z-index: "+strconv.Itoa(highestZIndex+1)+"; border: solid 2px #FF0000;")
      ghostWindow.Get("style").Set("left", Ftoa(startX) + "px")
      ghostWindow.Get("style").Set("top", Ftoa(startY) + "px")
      js.Global().Get("document").Get("body").Call("appendChild", ghostWindow)

      if verbose {Print("Resizing initiated with freeform selection.")}
    }
    return nil
  }))
}

