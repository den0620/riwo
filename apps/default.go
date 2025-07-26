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
	bg := wm.ThemeMap["green"]["faded"]
	mg := wm.ThemeMap["green"]["vivid"]
	fg := wm.ThemeMap["green"]["normal"]

	container := wm.Create()
	container.
		Style("display", "grid").
		Style("gridTemplateColumns", "repeat(auto-fit, minmax(120px, 1fr))").
		Style("background", bg).
		Style("gap", "5%").
		Style("padding", "5%").
		Style("height", "100%")

	title := wm.CreateKnown("h3")
	title.
		Inner("Applications").
		Style("gridColumn", "1 / -1").
		Style("fontSize", "24px").
		Style("color", fg).
		Style("textAlign", "center").
		Style("marginBotton", "20px").
		Mount(container)

	// This is an system Application
	// Get other registered applications
	for appName, appInit := range AppRegistry {
		buttonContainer := wm.
			Create().
			Style("textAlign", "center")

		appButton := wm.
			Create().
			Style("color", fg).
			Style("background", mg).
			Style("cursor", "url(assets/cursor-inverted.svg), auto").
			Style("padding", "15px").
			Style("borderRadius", "0").Style("border", "solid "+mg).
			Style("transition", "all 0.2s ease").
			Style("userSelect", "none")

		// prepare callbacks
		init := func(this js.Value, args []js.Value) interface{} {
			if wm.Verbose {
				wm.Print("App " + appName + " selected")
			}
			appInit(window)
			return nil
		}
		over := func(this js.Value, args []js.Value) interface{} {
			appButton.
				Style("background", fg).
				Style("color", bg)

			return nil
		}
		out := func(this js.Value, args []js.Value) interface{} {
			appButton.
				Style("background", bg).
				Style("color", fg)

			return nil
		}

		appButton.
			Inner(appName).
			Listen("mousedown", init).
			Listen("mouseover", over).
			Listen("mouseout", out)

		buttonContainer.
			Append(appButton).
			Mount(container)

	}

	window.Content.
		Inner("").
		Append(container)
}
