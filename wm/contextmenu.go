/*
Context Menu, its callbacks and listeners
*/

package wm

import (
	"strconv"
	"syscall/js"
)

// InitializeContextMenu creates and sets up the global context menu.
// It adds both default options and (if available) custom entries from the current window.
func InitializeContextMenu() {
	document := js.Global().Get("document")
	body := document.Get("body")

	// Create the context menu container.
	ContextMenu = document.Call("createElement", "div")
	ContextMenu.Set("id", "contextMenu")
	menuStyle := "position: absolute; display: none; background-color: " +
		GetColor["green"]["faded"] + "; border: solid " + GetColor["green"]["normal"] +
		"; padding: 0; text-align: center;"
	ContextMenu.Set("style", menuStyle)
	body.Call("appendChild", ContextMenu)

	// Pre-create default options.
	moveOption := CreateMenuOption("Move")
	newOption := CreateMenuOption("New")
	resizeOption := CreateMenuOption("Resize")
	deleteOption := CreateMenuOption("Delete")
	hideOption := CreateMenuOption("Hide")

	// Set up default options event handlers.
	moveOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsMovingMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			ContextMenu.Get("style").Set("display", "none")
			if Verbose {
				Print("Move mode activated.")
			}
		}
		return nil
	}))
	resizeOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsResizingMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			ContextMenu.Get("style").Set("display", "none")
			if Verbose {
				Print("Resize mode activated.")
			}
		}
		return nil
	}))
	deleteOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsDeleteMode = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			ContextMenu.Get("style").Set("display", "none")
			if Verbose {
				Print("Delete mode activated.")
			}
		}
		return nil
	}))
	hideOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsHiding = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			ContextMenu.Get("style").Set("display", "none")
			if Verbose {
				Print("Hide mode activated.")
			}
		}
		return nil
	}))

	// Background-specific options
	newOption.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsNewMode = true
			IsDragging = false
			StartX = 0
			StartY = 0
			js.Global().Get("document").Get("body").Get("style").Set("cursor", "url(assets/cursor-select.svg) 12 12, auto")
			ContextMenu.Get("style").Set("display", "none")
			if Verbose {
				Print("New mode activated. Select an area to create a window.")
			}
		}
		return nil
	}))

	// Omit browser's context menu
	body.Call("addEventListener", "contextmenu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		//justSelected = false
		return nil
	}))
	// And call ours
	body.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		if (args[0].Get("button").Int() == 2) && !JustSelected && !IsMovingMode && !IsResizingMode && !IsDeleteMode && !IsNewMode && !IsHiding {
			// Clear all previous menu items.
			for ContextMenu.Get("firstChild").Truthy() {
				ContextMenu.Call("removeChild", ContextMenu.Get("firstChild"))
			}
			// Add default options.
			ContextMenu.Call("appendChild", moveOption)
			ContextMenu.Call("appendChild", newOption)
			ContextMenu.Call("appendChild", resizeOption)
			ContextMenu.Call("appendChild", deleteOption)
			ContextMenu.Call("appendChild", hideOption)
			// Clear active window if deleted (or something)
			if (CurrentWindow != nil) && (CurrentWindow.ID == -1) {
				CurrentWindow = nil
			}
			// Append custom context menu entries from the active window, if any.
			if ActiveWindow.Truthy() && CurrentWindow != nil && CurrentWindow.ContextEntries != nil {
				for _, customOption := range CurrentWindow.ContextEntries {
					opt := CreateMenuOption(customOption.Name)
					opt.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true
						// Execute default callback for this option.
						customOption.Callback()
						ContextMenu.Get("style").Set("display", "none") // hide menu after click
						JustSelected = false
						if Verbose {
							Print("Custom option " + customOption.Name + " called")
						}
						return nil
					}))
					ContextMenu.Call("appendChild", opt)
				}
			}
			// Append hidden windows' unhides, if any.
			if len(ContextMenuHides) > 0 {
				for _, hiddenWindowButton := range ContextMenuHides {
					ContextMenu.Call("appendChild", hiddenWindowButton)
				}
			}

			// Position and show the menu.
			ContextMenu.Get("style").Set("z-index", strconv.Itoa(HighestZIndex+10))
			ContextMenu.Get("style").Set("left", strconv.Itoa(args[0].Get("clientX").Int())+"px")
			ContextMenu.Get("style").Set("top", strconv.Itoa(args[0].Get("clientY").Int())+"px")
			ContextMenu.Get("style").Set("display", "block")
		}
		JustSelected = false
		return nil
	}))
}

// CreateMenuOption creates a new context menu option element.
func CreateMenuOption(optionText string) js.Value {
	document := js.Global().Get("document")
	option := document.Call("createElement", "div")
	option.Set("innerText", optionText)
	option.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
	option.Get("style").Set("padding", "10px")
	option.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option.Get("style").Set("background-color", GetColor["green"]["vivid"])
		return nil
	}))
	option.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option.Get("style").Set("background-color", GetColor["green"]["faded"])
		return nil
	}))
	return option
}

// RemoveMenuOption removes a given menu option from the context menu.
func RemoveMenuOption(option js.Value) {
	document := js.Global().Get("document")
	menu := document.Call("getElementById", "contextMenu")
	menu.Call("removeChild", option)
}
