/*
App globals, primarily AppRegistry
*/

package apps

import (
	"riwo/wm"
)

// AppRegistry holds all available app functions.
// Each app should register itself (typically in its init function).
var AppRegistry = make(map[string]func(*wm.RiwoWindow))


// Simply apply theme to button
func applyThemeToButton(e *wm.RiwoObject, theme map[string]string) {
	e.
		Style("cursor", wm.CursorInvertUrl).
		Style("padding", "10px 20px").
		Style("color", "#000000").
		Style("backgroundColor", theme["faded"]).
		Style("border", "solid 2px "+theme["vivid"]).
		Style("borderRadius", "0")
}
