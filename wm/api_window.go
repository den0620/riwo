/*
Window struct for WM
and its functions + listeners
*/

package wm

import (
	"strconv"
	"syscall/js"
)

// Define a named type for the context entry.
type ContextEntry struct {
	Name     string
	Callback func()
}

// Type `RiwoWindow` manages single abstract window’s properties.
type RiwoWindow struct {
	ID          int            // For the most part unites DOM object and Go object
	Content     *RiwoElement   // Connected DOM element.
	MenuEntries []ContextEntry // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
}

// CreateWindow
// Creates a new Window, sets up its DOM element, and returns a pointer.
func CreateWindow(x, y, w, h, content string) *RiwoWindow {
	WindowCount++
	id := WindowCount

	document := js.Global().Get("document")
	body := document.Get("body")

	bodyContent := CreateFrom(&body)

	windowContent := Create().
		Style("overflow", "hidden").
		Style("position", "absolute").
		Style("width", w).
		Style("height", h).
		Style("top", y).
		Style("left", x).
		Style("z-index", strconv.Itoa(HighestZIndex)).
		Style("background-color", "#f0f0f0").
		Style("border", "solid #55AAAA").
		Style("padding", "0").
		Set("id", id).  // <-- assing shared ID
		Inner(content). // <-- spookie-dookie inner HTML
		Mount(bodyContent)

	// Logging
	if Verbose {
		Print("Generated window's ID (wid) is \"" +
			strconv.Itoa(id) + "\"")
	}

	window := &RiwoWindow{
		ID:      id,
		Content: windowContent,
		// No custom ContextEntries
	}

	CurrentWindow = window
	ActiveWindow = *windowContent

	AllWindows[strconv.Itoa(window.ID)] = window // <-- why string?????

	// Bring to front when clicked
	windowContent.Callback("mousedown", func(this js.Value, args []js.Value) interface{} {
		if !IsResizingInit {
			CurrentWindow = window
			ActiveWindow = *windowContent
		}

		// Right-click (RMB) on the window to select it for resizing, second right-click activates resizing
		if IsResizingMode && !IsResizingInit && args[0].Get("button").Int() == 2 {
			// First RMB hold - Select the window for resizing
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			if Verbose {
				Print("First right-click: Window selected for resizing.")
			}

			windowContent.Style("z-index", strconv.Itoa(HighestZIndex))
			HighestZIndex++
			IsResizingInit = true

			bodyContent.Style("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
		}

		// Mouse down event for selecting and dragging the window (click brings it to front)
		if !IsResizingInit {
			HighestZIndex++
			windowContent.Style("z-index", strconv.Itoa(HighestZIndex))
			if Verbose {
				Print("Window brought to front.")
			}

			if IsMovingMode && args[0].Get("button").Int() == 2 {
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				//JustSelected = true
				StartX = args[0].Get("clientX").Float() - windowContent.From("offsetLeft").Float()
				StartY = args[0].Get("clientY").Float() - windowContent.From("offsetTop").Float()
				IsDragging = true

				// Create ghost window
				rect := windowContent.Call("getBoundingClientRect")
				width := rect.Get("width").Float()
				height := rect.Get("height").Float()

				// Ensure ghost window is above everything during drag
				GhostWindow = *Create().
					Style("left", Ftoa(windowContent.From("offsetLeft").Float())+"px").
					Style("top", Ftoa(windowContent.From("offsetTop").Float())+"px").
					Style("position", "absolute").Style("z-index", strconv.Itoa(HighestZIndex+1)).
					Style("width", Ftoa(width)+"px").
					Style("height", Ftoa(height)+"px").
					Style("border", "solid 2px #FF0000").
					Style("cursor", "url(assets/cursor-drag.svg) 12 12, auto").
					Mount(bodyContent) // |<-- Append it to bodyContent

				JustSelected = true
				if Verbose {
					Print("Dragging initiated with ghost window.")
				}
			}

			if IsHiding && args[0].Get("button").Int() == 2 {
				// Hide window
				args[0].Call("preventDefault")
				args[0].Call("stopPropagation")
				JustSelected = true
				IsHiding = false

				hiddenWindowOption := CreateMenuOption("wid " + strconv.Itoa(window.ID))

				if windowContent.From("title").String() != "" {
					hiddenWindowOption = CreateMenuOption(windowContent.From("title").String())
				}

				hiddenWindowOption.Set("id", "menuopt"+strconv.Itoa(window.ID))

				// ??? option activation
				hiddenWindowOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					if args[0].Get("button").Int() == 2 {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true

						RemoveMenuOption(hiddenWindowOption)
						windowContent.Style("display", "block")

						ContextMenu.Get("style").Set("display", "none")

						// Delete by value
						for index, value := range ContextMenuHides {
							if value.Get("id").String() == hiddenWindowOption.Get("id").String() {
								ContextMenuHides = append(ContextMenuHides[:index], ContextMenuHides[index+1:]...)
							}
						}
						JustSelected = false
						if Verbose {
							Print("Unhide activated.")
						}
					}
					return nil
				}))
				ContextMenuHides = append(ContextMenuHides, hiddenWindowOption)
				windowContent.Style("display", "none")

				bodyContent.Style("cursor", "url(assets/cursor.svg), auto")

				JustSelected = false
				if Verbose {
					Print("WID " + strconv.Itoa(window.ID) + " hidden")
				}

			}
		}
		// Right-click (RMB) deletes the window in delete mode
		if IsDeleteMode && args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			RemoveWindow(window)
			IsDeleteMode = false

			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor.svg), auto")

			JustSelected = false
			if Verbose {
				Print("Window deleted.")
			}
		}
		return nil
	})

	return window
}

// Position
// Sets position and dimensions for the window. (Actually useless)
func (w *RiwoWindow) Position(newX, newY, newWidth, newHeight string) {
	w.Content.
		Style("left", newX).
		Style("top", newY).
		Style("width", newWidth).
		Style("height", newHeight)
}

// RemoveWindow
// Deletes the window from DOM and Go.
func RemoveWindow(w *RiwoWindow) {
	w.ID = -1                              // Remove reference for apps
	w.Content.Call("remove")               // Remove html part
	delete(AllWindows, strconv.Itoa(w.ID)) // Remove from list
}
