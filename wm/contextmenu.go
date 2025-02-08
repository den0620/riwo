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

