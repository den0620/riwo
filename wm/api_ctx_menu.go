/*
Context Menu, its callbacks and listeners
*/

package wm

import (
	"strconv"
	"syscall/js"
)

type MenuObject struct {
	dom js.Value
}

// DOM
// returns DOM data of current object
func (i *MenuObject) DOM() js.Value {
	return i.dom
}

func (i *MenuObject) Style(key string, value string) *MenuObject {
	i.dom.Get("style").Set(key, value)
	return i
}

func (i *MenuObject) Listen(event string, callback func(this js.Value, args []js.Value) interface{}) *MenuObject {
	i.dom.Call("addEventListener", event, js.FuncOf(callback))

	return i
}
func (i *MenuObject) Text(text string) *MenuObject {
	i.dom.Set("innerHTML", text)
	return i
}

// Depends on MenuItem and needed for MenuItem
// struct
func CreateMenuObject(text string) MenuObject {
	item := MenuObject{
		dom: js.Global().Get("document").Call("createElement", "div"),
	}
	item.
		Text(text).
		Style("cursor", CursorInvertUrl).
		Style("padding", "10px").
		Listen("mouseover", func(this js.Value, args []js.Value) interface{} {
			item.Style("background-color", ThemeMap["green"]["vivid"])
			return nil
		}).
		Listen("mouseout", func(this js.Value, args []js.Value) interface{} {
			item.Style("background-color", ThemeMap["green"]["faded"])
			return nil
		})

	return item
}

// InitializeContextMenu creates and sets up the global context menu.
// It adds both default options and (if available)
// custom entries from the current window.
func InitializeContextMenu() {
	document := js.Global().Get("document")
	b := document.Get("body")
	body := CreateFrom(&b)

	// Create the context menu container.
	ContextMenu = *Create().
		Set("id", "contextMenu").
		Style("position", "absolute").
		Style("display", "none").
		Style("background-color", ThemeMap["green"]["normal"]).
		Style("border", "solid "+ThemeMap["green"]["normal"]).
		Style("padding", "0").
		Style("text-align", "center")

	body.Append(&ContextMenu)

	// Pre-create default options.
	moveOption := CreateMenuObject("Move")
	newOption := CreateMenuObject("New")
	resizeOption := CreateMenuObject("Resize")
	deleteOption := CreateMenuObject("Delete")
	hideOption := CreateMenuObject("Hide")

	// Set up default options event handlers.
	moveOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			IsMovingMode = true

			js.Global().Get("document").Get("body").Get("style").Set("cursor", CursorSelectUrl)

			ContextMenu.Style("display", "none")
			if Verbose {
				Print("Move mode activated.")
			}
		}
		return nil
	})
	resizeOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			IsResizingMode = true

			js.Global().Get("document").Get("body").Get("style").Set("cursor", CursorSelectUrl)

			ContextMenu.Style("display", "none")
			if Verbose {
				Print("Resize mode activated.")
			}
		}
		return nil
	})
	deleteOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			IsDeleteMode = true

			js.Global().Get("document").Get("body").Get("style").Set("cursor", CursorSelectUrl)

			ContextMenu.Style("display", "none")
			if Verbose {
				Print("Delete mode activated.")
			}
		}
		return nil
	})
	hideOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")
			JustSelected = true
			IsHiding = true
			js.Global().Get("document").Get("body").Get("style").Set("cursor", CursorSelectUrl)
			ContextMenu.Style("display", "none")
			if Verbose {
				Print("Hide mode activated.")
			}
		}
		return nil
	})

	// Background-specific options
	newOption.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("button").Int() == 2 {
			args[0].Call("preventDefault")
			args[0].Call("stopPropagation")

			JustSelected = true
			IsNewMode = true
			IsDragging = false

			StartX = 0
			StartY = 0

			js.Global().Get("document").Get("body").Get("style").Set("cursor", CursorSelectUrl)
			ContextMenu.Style("display", "none")
			if Verbose {
				Print("New mode activated. Select an area to create a window.")
			}
		}
		return nil
	})

	// Omit browser's context menu
	body.Listen("contextmenu", func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		//justSelected = false
		return nil
	})

	// And call ours
	body.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")
		if (args[0].Get("button").Int() == 2) && !JustSelected && !IsMovingMode && !IsResizingMode && !IsDeleteMode && !IsNewMode && !IsHiding {
			// Clear all previous menu items.
			for ContextMenu.From("firstChild").Truthy() {
				ContextMenu.DOM().Call("removeChild", ContextMenu.From("firstChild"))
			}

			// Add default options. <-- problem
			ContextMenu.AppendByDom(
				moveOption.DOM(),
				newOption.DOM(),
				resizeOption.DOM(),
				deleteOption.DOM(),
				hideOption.DOM(),
			)

			// Clear active window if deleted (or something)
			if (CurrentWindow != nil) && (CurrentWindow.ID == -1) {
				CurrentWindow = nil
			}

			// Append custom context menu entries from the active window, if any.
			if ActiveWindow.DOM().Truthy() && CurrentWindow != nil && CurrentWindow.MenuEntries != nil {
				for _, customOption := range CurrentWindow.MenuEntries {
					opt := CreateMenuObject(customOption.Name)
					opt.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
						args[0].Call("preventDefault")
						args[0].Call("stopPropagation")
						JustSelected = true
						// Execute default callback for this option.
						customOption.Callback()
						ContextMenu.Style("display", "none") // hide menu after click
						JustSelected = false
						if Verbose {
							Print("Custom option " + customOption.Name + " called")
						}
						return nil
					})
					ContextMenu.AppendByDom(opt.DOM())
				}
			}
			// Append hidden windows' unhides, if any.
			if len(ContextMenuHides) > 0 {
				for _, hiddenWindowButton := range ContextMenuHides {
					ContextMenu.AppendByDom(hiddenWindowButton)
				}
			}

			// Position and show the menu.
			ContextMenu.
				Style("z-index", strconv.Itoa(HighestZIndex+10)).
				Style("left", strconv.Itoa(args[0].Get("clientX").Int())+"px").
				Style("top", strconv.Itoa(args[0].Get("clientY").Int())+"px").
				Style("display", "block")
		}
		JustSelected = false
		return nil
	})
}

// RemoveMenuOption
// removes a given menu option from the context menu.
func RemoveMenuOption(option js.Value) {
	document := js.Global().Get("document")
	menu := document.Call("getElementById", "contextMenu")
	menu.Call("removeChild", option)
}
