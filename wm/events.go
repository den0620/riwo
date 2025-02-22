/*
(Almost) All javascript DOM events here in only copy
for the sake of optimization
*/

package wm

import (
	"strconv"
	"syscall/js"
)

// InitializeGlobalMouseEvents sets up global mouse event listeners.
func InitializeGlobalMouseEvents() {
	// Moving mouse
	js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// In "Move" mode
		if IsDragging && IsMovingMode && GhostWindow.Truthy() {
			x := args[0].Get("clientX").Float() - StartX
			y := args[0].Get("clientY").Float() - StartY
			GhostWindow.Get("style").Set("left", Ftoa(x)+"px")
			GhostWindow.Get("style").Set("top", Ftoa(y)+"px")
			if Verbose {
				Print("Ghost window is moving.")
			}
		}

		// In "Resize" mode adjust ghost window size
		if GhostWindow.Truthy() && IsResizingMode && IsResizingInit && IsDragging && !IsMovingMode {
			currentX := args[0].Get("clientX").Float()
			currentY := args[0].Get("clientY").Float()
			// Calculate and update ghost window size and position based on selection
			width := currentX - StartX
			height := currentY - StartY
			GhostWindow.Get("style").Set("width", Ftoa(width)+"px")
			GhostWindow.Get("style").Set("height", Ftoa(height)+"px")
			// Handle direction to ensure ghost window moves according to selection direction
			if width < 0 {
				GhostWindow.Get("style").Set("left", Ftoa(currentX)+"px")
				GhostWindow.Get("style").Set("width", Ftoa(-width)+"px")
			}
			if height < 0 {
				GhostWindow.Get("style").Set("top", Ftoa(currentY)+"px")
				GhostWindow.Get("style").Set("height", Ftoa(-height)+"px")
			}
			if Verbose {
				Print("Ghost window is resizing with freeform selection.")
			}
		}

		// In "New" mode adjusting selection
		if GhostWindow.Truthy() && IsNewMode && IsDragging {
			currentX := args[0].Get("clientX").Float()
			currentY := args[0].Get("clientY").Float()
			// Calculate and update ghost window size and position based on selection
			width := currentX - StartX
			height := currentY - StartY
			GhostWindow.Get("style").Set("width", Ftoa(width)+"px")
			GhostWindow.Get("style").Set("height", Ftoa(height)+"px")
			// Handle direction to ensure ghost window moves according to selection direction
			if width < 0 {
				GhostWindow.Get("style").Set("left", Ftoa(currentX)+"px")
				GhostWindow.Get("style").Set("width", Ftoa(-width)+"px")
			}
			if height < 0 {
				GhostWindow.Get("style").Set("top", Ftoa(currentY)+"px")
				GhostWindow.Get("style").Set("height", Ftoa(-height)+"px")
			}
			if Verbose {
				Print("New window selection resizing.")
			}
		}

		return nil
	}))

	// Mouse down
	js.Global().Get("document").Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// In "Resize" Second RMB after choosing window
		// Also "New"'s initial area selecting
		if (IsResizingMode || IsNewMode) && args[0].Get("button").Int() == 2 && (IsResizingInit || IsNewMode) {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			// Second RMB hold - Start resizing by creating a selection anywhere
			if Verbose {
				Print("Second right-click: Resizing initiated.")
			}
			IsDragging = true
			// Reset selection start position
			StartX = args[0].Get("clientX").Float()
			StartY = args[0].Get("clientY").Float()
			// Create ghost window for resizing
			GhostWindow = js.Global().Get("document").Call("createElement", "div")
			GhostWindow.Set("style", "position: absolute; z-index: "+strconv.Itoa(HighestZIndex+1)+"; border: solid 2px #FF0000;")
			GhostWindow.Get("style").Set("left", Ftoa(StartX)+"px")
			GhostWindow.Get("style").Set("top", Ftoa(StartY)+"px")
			js.Global().Get("document").Get("body").Call("appendChild", GhostWindow)
			if Verbose {
				Print("Resizing initiated with freeform selection.")
			}
		}
		// Cancel menu and reset all modes on leftclick.
		if args[0].Get("button").Int() == 0 {
			if ContextMenu.Type() == js.TypeUndefined {
				Print("ContextMenu is UNDEFINED!")
				// Or handle the error appropriately
			} else {
				ContextMenu.Get("style").Set("display", "none")
			}
			if GhostWindow.Truthy() && IsDragging {
				GhostWindow.Call("remove")
				GhostWindow = js.Null()
			}
			IsDragging = false
			IsMovingMode = false
			IsResizingMode = false
			IsResizingInit = false
			JustSelected = false
			IsDeleteMode = false
			IsNewMode = false
			IsHiding = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
		}
		return nil
	}))

	// Mouse up
	js.Global().Get("document").Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// In "move" stop dragging (Teleport to ghost)
		if IsMovingMode && IsDragging {
			args[0].Call("stopPropagation")
			IsDragging = false
			IsMovingMode = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
			// Move the window to the ghost's position
			if GhostWindow.Truthy() && ActiveWindow.Truthy() {
				ActiveWindow.Get("style").Set("left", GhostWindow.Get("style").Get("left"))
				ActiveWindow.Get("style").Set("top", GhostWindow.Get("style").Get("top"))
				GhostWindow.Call("remove") // Remove ghost window
				GhostWindow = js.Null()    // Reset ghost window reference
			}
			JustSelected = false
			if Verbose {
				Print("Dragging ended and window teleported to ghost position.")
			}
		}

		// In "Resize" stop selecting (Teleport to ghost)
		if GhostWindow.Truthy() && ActiveWindow.Truthy() && IsResizingMode && IsResizingInit && IsDragging && !IsMovingMode {
			args[0].Call("stopPropagation")
			IsResizingMode = false
			IsResizingInit = false
			IsDragging = false
			// Replace all dimensions with ghost's ones
			ActiveWindow.Get("style").Set("left", GhostWindow.Get("style").Get("left"))
			ActiveWindow.Get("style").Set("top", GhostWindow.Get("style").Get("top"))
			ActiveWindow.Get("style").Set("width", GhostWindow.Get("style").Get("width"))
			ActiveWindow.Get("style").Set("height", GhostWindow.Get("style").Get("height"))
			GhostWindow.Call("remove") // Remove the ghost window
			GhostWindow = js.Null()    // Reset ghost window reference
			// Reset cursor
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
			JustSelected = false
			if Verbose {
				Print("Resizing completed and window resized to match selection.")
			}
		}

		// In "New" make new window
		if IsNewMode && IsDragging {
			args[0].Call("stopPropagation")
			IsNewMode = false
			IsDragging = false
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")
			// Create a new window at the ghost window's position and size
			if GhostWindow.Truthy() {
				x := GhostWindow.Get("style").Get("left").String()
				y := GhostWindow.Get("style").Get("top").String()
				width := GhostWindow.Get("style").Get("width").String()
				height := GhostWindow.Get("style").Get("height").String()
				GhostWindow.Call("remove")
				GhostWindow = js.Null()
				neuwindow := WindowCreate(x, y, width, height, "")
				JustSelected = false
				if Verbose {
					Print("New window created at selected area.")
				}
				// Intends existance of default app
				APP_default := js.Global().Get("LaunchDefault")
				go APP_default.Invoke(neuwindow.ID)
				if Verbose {
					Print("LaunchDefault attached to the new window")
				}
			}
		}

		return nil
	}))
	return
}
