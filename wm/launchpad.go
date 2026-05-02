package wm

import (
	"syscall/js"
)

func guestAppLabels() []string {
	k := js.Global().Get("__riwoKernel")
	if !k.Truthy() {
		return nil
	}
	arr := k.Call("listGuestApps")
	if arr.Type() == js.TypeNull || arr.Type() == js.TypeUndefined {
		return nil
	}
	n := arr.Get("length").Int()
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		v := arr.Index(i).String()
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// attachExitMenu registers a wm-local context menu row to leave a guest and reopen the launcher grid
func attachExitMenu(w *RiwoWindow) {
	if w == nil {
		return
	}
	w.MenuEntries = []ContextEntry{
		{Name: "Exit", Callback: func() { ReturnToLaunchpad(w) }},
	}
}

// ReturnToLaunchpad stops the guest (kernel bookkeeping), clears its menus, and rebuilds this window's launcher
func ReturnToLaunchpad(w *RiwoWindow) {
	if w == nil {
		return
	}
	js.Global().Get("__riwoKernel").Call("disposeGuestForWindow", w.ID)
	OpenLaunchpad(w)
}

// OpenLaunchpad fills the window with app tiles; each tile asks the kernel to load a wasm guest
func OpenLaunchpad(window *RiwoWindow) {
	window.Title = "Launchpad"
	window.MenuEntries = nil

	bg := GetBackgroundColorStr("green")
	mg := GetBorderColorStr("green")
	fg := GetFontColorStr("green")

	container := Create()
	container.
		Style("display", "grid").
		Style("gridTemplateColumns", "repeat(auto-fit, minmax(6rem, 1fr))").
		Style("background", bg).
		Style("gap", "0.25rem").
		Style("padding", "0.25rem").
		Style("height", "100%")

	appNames := guestAppLabels()
	for _, appName := range appNames {
		btnApp := appName
		buttonContainer := Create().
			Style("textAlign", "center")

		appButton := Create().
			Style("color", "#000000").
			Style("background", bg).
			Style("cursor", CursorInvertUrl).
			Style("padding", "1rem").
			Style("width", "auto").
			Style("height", "auto").
			Style("borderRadius", "0").
			Style("border", "solid "+mg).
			Style("userSelect", "none")

		init := func(this js.Value, args []js.Value) interface{} {
			JSLog("App " + btnApp + " selected")
			window.Title = btnApp
			attachExitMenu(window)
			KernelSpawnGuest(btnApp, window)
			return nil
		}
		out := func(this js.Value, args []js.Value) interface{} {
			appButton.Style("background", bg)
			return nil
		}
		over := func(this js.Value, args []js.Value) interface{} {
			appButton.Style("background", fg)
			return nil
		}

		appButton.
			Inner(btnApp).
			Listen("mousedown", init).
			Listen("mouseout", out).
			Listen("mouseover", over)

		buttonContainer.
			Append(appButton).
			Mount(container)
	}

	window.Content.
		Inner("").
		Append(container)
}
