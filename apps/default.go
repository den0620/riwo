package apps

import (
	"riwo/wm"
	"syscall/js"
)

func init() {
	// Register the default app itself.
	AppRegistry["Default"] = Construct
}

func Construct(window *wm.RiwoWindow) {
	window.Title = "Launchpad"

	bg := wm.GetBackgroundColorStr("green")
	mg := wm.GetBorderColorStr("green")
	fg := wm.GetFontColorStr("green")

	container := wm.Create()
	container.
		Style("display", "grid").
		Style("gridTemplateColumns", "repeat(auto-fit, minmax(6rem, 1fr))").
		Style("background", bg).
		Style("gap", "0.25rem").
		Style("padding", "0.25rem").
		Style("height", "100%")

	/*
		title := wm.Create()
		title.
			Inner("Applications").
			Style("gridColumn", "1 / -1").
			Style("fontSize", "24px").
			Style("color", mg).
			Style("textAlign", "center").
			Style("margin", "20px").
			Mount(container)
	*/

	// This is an system Application
	// Get other registered applications
	for appName, appInit := range AppRegistry {
		buttonContainer := wm.
			Create().
			Style("textAlign", "center")

		appButton := wm.
			Create().
			Style("color", "#000000").
			Style("background", bg).
			Style("cursor", wm.CursorInvertUrl).
			Style("padding", "1rem").
			Style("width", "auto").
			Style("height", "auto").
			Style("borderRadius", "0").Style("border", "solid "+mg).
			Style("userSelect", "none")

		// prepare callbacks
		init := func(this js.Value, args []js.Value) interface{} {
			wm.JSLog("App " + appName + " selected")

			// warn: After window initizliation, it would be better if
			// known application name was applied to window title.
			//
			// Elsewhere we have replaced window content and previous app name (Launchpad #wid)
			window.Title = appName

			appInit(window)
			return nil
		}

		out := func(this js.Value, args []js.Value) interface{} {
			appButton.
				Style("background", bg)

			return nil
		}
		over := func(this js.Value, args []js.Value) interface{} {
			appButton.
				Style("background", fg)

			return nil
		}

		appButton.
			Inner(appName).
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
