package apps

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func init() {
	AppRegistry["Player"] = APP_dplayer
}

func APP_dplayer(window *wm.Window) {
	document := js.Global().Get("document")
	currentTheme := "blue"
	colorFG := wm.ThemeMap[currentTheme]["vivid"]
	colorBG := wm.ThemeMap[currentTheme]["faded"]
	colorMG := wm.ThemeMap[currentTheme]["normal"]

	// Create main container
	container := document.Call("createElement", "div")
	container.Get("style").Set("height", "100%")
	container.Get("style").Set("display", "flex")
	container.Get("style").Set("flexDirection", "column")
	container.Get("style").Set("justifyContent", "center")
	container.Get("style").Set("alignItems", "center")
	container.Get("style").Set("backgroundColor", wm.ThemeMap[currentTheme]["faded"])

	// Create audio element
	audio := document.Call("createElement", "audio")
	audio.Set("src", "about:blank") // Default empty source

	// Time display
	timeDisplay := document.Call("createElement", "div")
	timeDisplay.Set("textContent", "00:00 / 00:00")
	timeDisplay.Get("style").Set("color", wm.ThemeMap[currentTheme]["vivid"])
	//timeDisplay.Get("style").Set("fontFamily", "monospace")
	timeDisplay.Get("style").Set("textAlign", "center")
	timeDisplay.Get("style").Set("marginBottom", "10px")

	// Control buttons container
	controlsContainer := document.Call("createElement", "div")
	controlsContainer.Get("style").Set("display", "flex")
	controlsContainer.Get("style").Set("alignItems", "center")

	// Control buttons container
	controlsContainer2 := document.Call("createElement", "div")
	controlsContainer2.Get("style").Set("display", "flex")
	controlsContainer2.Get("style").Set("alignItems", "center")

	// Style for buttons
	styleButton := func(btn js.Value) {
		btn.Get("style").Set("padding", "5px 15px")
		btn.Get("style").Set("backgroundColor", colorBG)
		btn.Get("style").Set("color", "black")
		btn.Get("style").Set("border", "solid "+colorMG)
		btn.Get("style").Set("borderRadius", "0")
		btn.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
		btn.Get("style").Set("font-family", "monospace")
		btn.Get("style").Set("font-weight", "bold")
		// Hover effects
		btn.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			btn.Get("style").Set("background", colorFG)
			btn.Get("style").Set("color", colorBG)
			return nil
		}))
		btn.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			btn.Get("style").Set("background", colorBG)
			btn.Get("style").Set("color", "black")
			return nil
		}))
	}

	// Control buttons
	prevBtn := document.Call("createElement", "div")
	prevBtn.Set("textContent", "[<<]")
	styleButton(prevBtn)

	playBtn := document.Call("createElement", "div")
	playBtn.Set("textContent", "[I>]")
	styleButton(playBtn)

	nextBtn := document.Call("createElement", "div")
	nextBtn.Set("textContent", "[>>]")
	styleButton(nextBtn)

	// Volume controls
	volumeDown := document.Call("createElement", "div")
	volumeDown.Set("textContent", "[-]")
	styleButton(volumeDown)

	volumeUp := document.Call("createElement", "div")
	volumeUp.Set("textContent", "[+]")
	styleButton(volumeUp)

	// File selection button
	fileBtn := document.Call("createElement", "div")
	fileBtn.Set("textContent", "File")
	styleButton(fileBtn)

	// Hidden file input
	fileInput := document.Call("createElement", "input")
	fileInput.Set("type", "file")
	fileInput.Set("accept", "*/*")
	fileInput.Get("style").Set("display", "none")

	// Status text
	statusText := document.Call("createElement", "div")
	statusText.Set("textContent", "No file loaded")
	statusText.Get("style").Set("color", wm.ThemeMap[currentTheme]["vivid"])
	statusText.Get("style").Set("margin", "10px")

	// Format time helper
	formatTime2Digits := func(n int) string {
		if n < 10 {
			return "0" + strconv.Itoa(n)
		}
		return strconv.Itoa(n)
	}
	formatTime := func(seconds float64) string {
		min := int(seconds) / 60
		sec := int(seconds) % 60
		return strconv.Itoa(min) + ":" + formatTime2Digits(sec)
	}

	// Skip forward/backward by 10 seconds
	prevBtn.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.Get("currentTime").Float()
		audio.Set("currentTime", currentTime-10)
		return nil
	}))

	nextBtn.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.Get("currentTime").Float()
		audio.Set("currentTime", currentTime+10)
		return nil
	}))

	// Event handlers
	playBtn.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if audio.Get("paused").Bool() {
			audio.Call("play")
			playBtn.Set("textContent", "[||]")
		} else {
			audio.Call("pause")
			playBtn.Set("textContent", "[I>]")
		}
		return nil
	}))

	volumeDown.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentVol := audio.Get("volume").Float()
		if currentVol > 0.1 {
			audio.Set("volume", currentVol-0.1)
		}
		return nil
	}))

	volumeUp.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentVol := audio.Get("volume").Float()
		if currentVol < 1.0 {
			audio.Set("volume", currentVol+0.1)
		}
		return nil
	}))

	fileBtn.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fileInput.Call("click")
		return nil
	}))

	// Audio time update handler
	audio.Call("addEventListener", "timeupdate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.Get("currentTime").Float()
		duration := audio.Get("duration").Float()
		if !js.Global().Call("isNaN", duration).Bool() {
			timeDisplay.Set("textContent", formatTime(currentTime)+" / "+formatTime(duration))
		}
		if currentTime == duration {
			audio.Call("pause")
			playBtn.Set("textContent", "[I>]")
		}
		return nil
	}))

	// File input change handler
	fileInput.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		files := fileInput.Get("files")
		if files.Length() > 0 {
			file := files.Index(0)
			url := js.Global().Get("URL").Call("createObjectURL", file)
			audio.Set("src", url)
			statusText.Set("textContent", ":: "+file.Get("name").String())
			playBtn.Set("textContent", "[I>]")
		}
		return nil
	}))

	// Add elements to controls container (upper)
	controlsContainer.Call("appendChild", prevBtn)
	controlsContainer.Call("appendChild", playBtn)
	controlsContainer.Call("appendChild", nextBtn)
	// Add elements to controls container (lower)
	controlsContainer2.Call("appendChild", volumeDown)
	controlsContainer2.Call("appendChild", volumeUp)
	controlsContainer2.Call("appendChild", fileBtn)

	// Add elements to main container
	container.Call("appendChild", timeDisplay)
	container.Call("appendChild", controlsContainer)
	container.Call("appendChild", statusText)
	container.Call("appendChild", controlsContainer2)
	container.Call("appendChild", fileInput)
	container.Call("appendChild", audio)

	// Clear window and add container
	window.DOM.Set("innerHTML", "")
	window.DOM.Call("appendChild", container)
}
