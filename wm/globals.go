/*
Global vars for WM
and common utils
*/

package wm

import (
	"strconv"
	"syscall/js"
)

var (
	ContextMenu      js.Value
	ContextMenuHides []js.Value
	IsDragging       bool
	IsMovingMode     bool
	IsResizingMode   bool
	IsResizingInit   bool
	JustSelected     bool
	IsDeleteMode     bool
	IsNewMode        bool
	IsHiding         bool
	StartX, StartY   float64
	CurrentWindow    *RiwoWindow            // Active Go Window
	ActiveWindow     js.Value               // Active JS window
	AllWindows       map[string]*RiwoWindow // All Go Windows
	GhostWindow      js.Value
	WindowCount      int       // Counter for creating multiple windows with unique z-index
	HighestZIndex    int  = 10 // Track the highest z-index for bringing windows to front
	Verbose          bool = false
	ThemeMap              = map[string]map[string]string{
		"monochrome": {
			"faded":  "#ffffff",
			"normal": "#777777",
			"vivid":  "#000000",
		},
		"red": {
			"faded":  "#ffeaea",
			"normal": "#df9595",
			"vivid":  "#bb5d5d",
		},
		"green": {
			"faded":  "#eaffea",
			"normal": "#88cc88",
			"vivid":  "#448844",
		},
		"blue": {
			"faded":  "#c0eaff",
			"normal": "#00aaff",
			"vivid":  "#0088cc",
		},
		"yellow": {
			"faded":  "#ffffea",
			"normal": "#eeee9e",
			"vivid":  "#99994c",
		},
		"aqua": {
			"faded":  "#eaffff",
			"normal": "#9eeeee",
			"vivid":  "#8888cc",
		},
		"gray": {
			"faded":  "#eeeeee",
			"normal": "#cccccc",
			"vivid":  "#888888",
		},
	}
)

func Print(value string) {
	js.Global().Get("console").Call("log", value)
}
func Ftoa(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}
