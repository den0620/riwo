/*
(Almost) All javascript DOM events here in only copy
for the sake of optimization
*/

package wm

import (
	"strconv"
	"syscall/js"
)

// InitializeGlobalMouseEvents
// sets up global mouse event listeners.
func InitializeGlobalMouseEvents() {
	// Moving mouse
	js.Global().Get("document").Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// In "Move" mode
		if IsDragging && IsMovingMode && GhostWindow.DOM().Truthy() {
			x := args[0].Get("clientX").Float() - StartX
			y := args[0].Get("clientY").Float() - StartY

			GhostWindow.
				Style("left", Ftoa(x)+"px").
				Style("top", Ftoa(y)+"px")

			if Verbose {
				Print("Ghost window is moving.")
			}
		}

		// In "Resize" mode adjust ghost window size
		if GhostWindow.DOM().Truthy() && IsResizingMode && IsResizingInit && IsDragging && !IsMovingMode {
			currentX := args[0].Get("clientX").Float()
			currentY := args[0].Get("clientY").Float()
			// Calculate and update ghost window size and position based on selection
			width := currentX - StartX
			height := currentY - StartY

			GhostWindow.
				Style("width", Ftoa(width)+"px").
				Style("height", Ftoa(height)+"px")

			// Handle direction to ensure ghost window moves according to selection direction
			if width < 0 {
				GhostWindow.
					Style("left", Ftoa(currentX)+"px").
					Style("width", Ftoa(-width)+"px")
			}
			if height < 0 {
				GhostWindow.
					Style("top", Ftoa(currentY)+"px").
					Style("height", Ftoa(-height)+"px")
			}
			if Verbose {
				Print("Ghost window is resizing with freeform selection.")
			}
		}

		// In "New" mode adjusting selection
		if GhostWindow.DOM().Truthy() && IsNewMode && IsDragging {
			currentX := args[0].Get("clientX").Float()
			currentY := args[0].Get("clientY").Float()

			// Calculate and update ghost window size and position based on selection
			width := currentX - StartX
			height := currentY - StartY
			GhostWindow.
				Style("width", Ftoa(width)+"px").
				Style("height", Ftoa(height)+"px")
			// Handle direction to ensure ghost window moves according to selection direction
			if width < 0 {
				GhostWindow.
					Style("left", Ftoa(currentX)+"px").
					Style("width", Ftoa(-width)+"px")
			}
			if height < 0 {
				GhostWindow.
					Style("top", Ftoa(currentY)+"px").
					Style("height", Ftoa(-height)+"px")
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
			GhostWindow = *Create().
				Style("position", "absolute").
				Style("z-index", strconv.Itoa(HighestZIndex+1)).
				Style("border", "solid 2px #FF0000")

			GhostWindow.
				Style("left", Ftoa(StartX)+"px").
				Style("top", Ftoa(StartY)+"px")

			js.Global().Get("document").Get("body").Call("appendChild", GhostWindow.DOM()) // <-- !!!

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
			if GhostWindow.DOM().Truthy() && IsDragging {
				GhostWindow.Call("remove")
				GhostWindow.Delete() // <-- erase DOM
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
			if GhostWindow.DOM().Truthy() && ActiveWindow.DOM().Truthy() {
				ActiveWindow.
					Style("left", GhostWindow.From("style").Get("left")).
					Style("top", GhostWindow.From("style").Get("top"))

				GhostWindow.Call("remove") // Remove ghost window
				GhostWindow.Delete()       // Reset ghost window reference
			}
			JustSelected = false
			if Verbose {
				Print("Dragging ended and window teleported to ghost position.")
			}
		}

		// In "Resize" stop selecting (Teleport to ghost)
		if GhostWindow.DOM().Truthy() && ActiveWindow.DOM().Truthy() && IsResizingMode && IsResizingInit && IsDragging && !IsMovingMode {
			args[0].Call("stopPropagation")
			IsResizingMode = false
			IsResizingInit = false
			IsDragging = false

			// Replace all dimensions with ghost's ones
			ActiveWindow.
				Style("left", GhostWindow.From("style").Get("left")).
				Style("top", GhostWindow.From("style").Get("top")).
				Style("width", GhostWindow.From("style").Get("width")).
				Style("height", GhostWindow.From("style").Get("height"))

			GhostWindow.Call("remove") // Remove the ghost window
			GhostWindow.Delete()       // Reset ghost window reference

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
			if GhostWindow.DOM().Truthy() {
				x := GhostWindow.From("style").Get("left").String()
				y := GhostWindow.From("style").Get("top").String()
				width := GhostWindow.From("style").Get("width").String()
				height := GhostWindow.From("style").Get("height").String()

				GhostWindow.Call("remove")
				GhostWindow.Delete()

				newWindow := CreateWindow(x, y, width, height, "")

				JustSelected = false
				if Verbose {
					Print("New window created at selected area.")
				}
				// Intends existance of default app
				defaultInit := js.Global().Get("LaunchDefault")
				go defaultInit.Invoke(newWindow.ID)

				if Verbose {
					Print("LaunchDefault attached to the new window")
				}
			}
		}

		return nil
	}))
}
