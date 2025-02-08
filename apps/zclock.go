package apps


import (
	"strconv"
	"syscall/js"
	"riwo/wm"
)


var (
	isSettings        bool = false  // Interface State
	localeUTC          int = 5      // Urals because I can
	colorTheme      string = "aqua" // Faded for background, vivid for foreground
)


func SwitchToSettings(window *Window) {
	// do shit
	colorBG = GetColor[colorTheme]["faded"]
	colorFG = GetColor[colorTheme]["vivid"]
	// Add here wheel for localeUTC
	// Add here list menu for colorTheme
}

func SwitchToFace(window *Window) {
	// do shit
	colorBG = GetColor[colorTheme]["faded"]
	colorFG = GetColor[colorTheme]["vivid"]
	// Make some kind of clock
}

func main(window *Window) {
	// Initialize custom context menu entry
	optionSettings js.Value = CreateMenuOption("settings")
	// Callback
	optionSettings.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if isSettings {
			SwitchToFace(*window)
		} else {
			SwitchToSettings(*window)
		}
		isSettings = !isSettings
		if Verbose {wm.Print("zclock \"settings\" activated, went to "+isSettings.String())}
	}))
	
}
