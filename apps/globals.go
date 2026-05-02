/*
Shared helpers for guest apps built as separate wasm binaries.
*/

package apps

import "riwo/wm"

func applyThemeToButton(e *wm.RiwoObject, theme map[string]string) {
	e.
		Style("cursor", wm.CursorInvertUrl).
		//Style("padding", "10px 20px").
		Style("color", "#000000").
		Style("backgroundColor", theme["faded"]).
		Style("border", "solid 2px "+theme["vivid"]).
		Style("borderRadius", "0")
}
