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
	ContextMenu      RiwoObject
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
	AllWindows       map[string]*RiwoWindow // All Go Windows
	CurrentWindow    *RiwoWindow            // Active Go Window
	ActiveWindow     RiwoObject             // Active JS window
	GhostWindow      RiwoObject
	WindowCount      int       // Counter for creating multiple windows with unique z-index
	HighestZIndex    int  = 10 // Track the highest z-index for bringing windows to front
	Verbose          bool = false
)

// / JSLog prints the string in JS console (if flag verbose was enabled)
func JSLog(value string) {
	if Verbose {
		js.
			Global().
			Get("console").
			Call("log", value)
	}
}
func Ftoa(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}
