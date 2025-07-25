package apps

import (
	"riwo/wm"
)

func init() {
	AppRegistry["testapp"] = AppTest
}

func AppTest(window *wm.Window) {
	// init element
	container := wm.Create()

	// prepare styles
	container.
		Style("background", wm.ThemeMap["green"]["faded"]).
		Style("gap", "5%").
		Style("padding", "5%").
		Style("height", "100%").
		Style("display", "grid")

	title := wm.Create()
	title.
		Inner("RIWO App"). // <-- text inside title container
		Style("fontSize", "24px").
		Style("textAlign", "center").
		Style("marginBotton", "20px").
		Mount(container) // <-- add element to parent

	window.DOM.Set("innerHTML", "")
	window.DOM.Call("appendChild", container.DOM()) // <-- ideas??
}
