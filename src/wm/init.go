package wm

import (
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
  isNewMode      bool
  isHiding       bool
  startX, startY float64
  activeWindow   js.Value
  ghostWindow    js.Value
  windowCount    int    // Counter for creating multiple windows with unique z-index
  highestZIndex  int = 10 // Track the highest z-index for bringing windows to front
  verbose        bool = false
)


func Print(value string) {
  js.Global().Get("console").Call("log", value)
}
func Ftoa(value float64) string {
  return strconv.FormatFloat(value, 'f', 6, 64)
}
func GoVerbose(this js.Value, args []js.Value) interface{} {
  verbose = !verbose
  Print("Verbose : " + strconv.FormatBool(verbose))
  return nil
}


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

func CreateDraggableWindow(x string, y string, width string, height string) interface{} {
  document := js.Global().Get("document")
  body := document.Get("body")

  windowCount++ // Increment window counter for unique z-index

  // Create the window element
  window := document.Call("createElement", "div")
  window.Set("style", "overflow: hidden; position: absolute; z-index: "+strconv.Itoa(highestZIndex)+"; left: "+x+"; top: "+y+"; width: "+width+"; height: "+height+"; background-color: #f0f0f0; border: solid #55AAAA; padding: 0;")
  window.Set("innerHTML", "<h3>Draggable Window "+strconv.Itoa(windowCount)+"</h3><p>html p</p>")
  window.Set("title", "Test" + strconv.Itoa(windowCount))
  window.Set("wid", strconv.Itoa(windowCount))
  window.Set("id", strconv.Itoa(windowCount))

  if verbose {Print("Generated window's title is \""+window.Get("title").String()+"\"; Window's ID (wid) is \""+window.Get("wid").String()+"\"")}

  body.Call("appendChild", window)

  // Prevent the context menu from opening on right-click for windows
  window.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    //args[0].Call("stopPropagation")
    if verbose {Print("Caught click on window")}
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
      if verbose {Print("First right-click: Window selected for resizing.")}
      highestZIndex++
      window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
      activeWindow = window
      isResizingInit = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
    }

    // Mouse down event for selecting and dragging the window (Left-click brings it to front)
    if !isResizingInit {
      // Left-click to bring window to front
      highestZIndex++
      window.Get("style").Set("z-index", strconv.Itoa(highestZIndex))
      if verbose {Print("Window brought to front.")}

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
        ghostWindow.Set("style", "position: absolute; z-index: "+strconv.Itoa(highestZIndex+1)+"; width: "+Ftoa(width)+"px; height: "+Ftoa(height)+"px; border: solid 2px #FF0000; cursor: url(assets/cursor-drag.svg) 12 12, auto;")
        ghostWindow.Get("style").Set("left", Ftoa(window.Get("offsetLeft").Float())+"px")
        ghostWindow.Get("style").Set("top", Ftoa(window.Get("offsetTop").Float())+"px")
        body.Call("appendChild", ghostWindow)

        if verbose {Print("Dragging initiated with ghost window.")}
      }
      if isHiding && args[0].Get("button").Int() == 2 {
        // Hide window
        args[0].Call("preventDefault")
        args[0].Call("stopPropagation")
        isHiding = false
        menu := document.Call("getElementById", "contextMenu")

	hidenWindowOption := CreateMenuOption(window.Get("title").String())
	hidenWindowOption.Set("id", "menuopt"+window.Get("wid").String())

        // Unhide option activation
        hidenWindowOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
          if args[0].Get("button").Int() == 2 {
            args[0].Call("preventDefault")
            args[0].Call("stopPropagation")
	    RemoveMenuOption(hidenWindowOption)
            window.Get("style").Set("display", "block")
            menu.Get("style").Set("display", "none")
            if verbose {Print("Unhide activated.")}
          }
          return nil
        }))
        menu.Call("appendChild", hidenWindowOption)
        if verbose {Print("option added\n")}

        window.Get("style").Set("display", "none")
	justSelected = false
        js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
        if verbose {Print("WID "+window.Get("wid").String()+" hidden")}
      }
    }

    // Right-click (RMB) deletes the window in delete mode
    if isDeleteMode && args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      window.Call("remove") // Delete the window
      isDeleteMode = false
      justSelected = false
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
      if verbose {Print("Window deleted.")}
    }
    return nil
  }))
    
  return window
}


func InitializeContextMenu() {
  document := js.Global().Get("document")
  body := document.Get("body")

  // Create the context menu with higher z-index
  menu := document.Call("createElement", "div")
  menu.Set("id", "contextMenu")
  menu.Set("style", "position: absolute; display: none; background-color: #EEFFEE; border: solid #8BCE8B; padding: 0; text-align: center;")
  body.Call("appendChild", menu)

  // Move, New, Resize, and Delete options
  moveOption := CreateMenuOption("Move")
  newOption := CreateMenuOption("New")
  resizeOption := CreateMenuOption("Resize")
  deleteOption := CreateMenuOption("Delete")
  hideOption := CreateMenuOption("Hide")

  menu.Call("appendChild", moveOption)
  menu.Call("appendChild", newOption)
  menu.Call("appendChild", resizeOption)
  menu.Call("appendChild", deleteOption)
  menu.Call("appendChild", hideOption)

  // Cancel menu and actions on left-click
  body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 0 {
      menu.Get("style").Set("display", "none")
      if ghostWindow.Truthy() && isDragging {ghostWindow.Call("remove")}
      isDragging = false
      isMovingMode = false
      isResizingMode = false
      isResizingInit = false
      justSelected = false // test
      isDeleteMode = false
      isNewMode = false
      isHiding = false
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
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
      menu.Get("style").Set("display", "none")
      if verbose {Print("Move mode activated.")}
    }
    return nil
  }))

  // New window activation
  newOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      isNewMode = true
      isDragging = false
      startX = 0
      startY = 0

      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
      menu.Get("style").Set("display", "none")
      if verbose {Print("New mode activated. Select an area to create a window.")}
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
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
      menu.Get("style").Set("display", "none")
      if verbose {Print("Resize mode activated.")}
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
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
      menu.Get("style").Set("display", "none")
      if verbose {Print("Delete mode activated.")}
    }
    return nil
  }))

  // Hide mode activation
  hideOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if args[0].Get("button").Int() == 2 {
      args[0].Call("preventDefault")
      args[0].Call("stopPropagation")
      justSelected = true
      isHiding = true
      js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
      menu.Get("style").Set("display", "none")
      if verbose {Print("Hide mode activated.")}
    }
    return nil
  }))

  // Global context menu activation
  body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    justSelected = false
    return nil
  }))

  body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    args[0].Call("preventDefault")
    if (args[0].Get("button").Int() == 2) && !justSelected && !isMovingMode && !isResizingMode && !isDeleteMode && !isNewMode && !isHiding {
      // Adjust z-index dynamically based on highestZIndex
      menu.Get("style").Set("z-index", strconv.Itoa(highestZIndex+10))
      menu.Get("style").Set("left", strconv.Itoa(args[0].Get("clientX").Int()) + "px")
      menu.Get("style").Set("top", strconv.Itoa(args[0].Get("clientY").Int()) + "px")
      menu.Get("style").Set("display", "block")
    }
    justSelected = false
    return nil
  }))
}

// Menu option creation
func CreateMenuOption(optionText string) js.Value {
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

// Menu option deletion
func RemoveMenuOption(option js.Value) {
  document := js.Global().Get("document")
  menu := document.Call("getElementById", "contextMenu")
  menu.Call("removeChild", option)
  return
}

