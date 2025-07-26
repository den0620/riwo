package apps

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func init() {
	AppRegistry["Player"] = PlayerConstruct
}

func APP_dplayer(window *wm.RiwoWindow) {
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
	timeDisplay.Get("style").Set("textAlign", "center")
	timeDisplay.Get("style").Set("marginBottom", "10px")
	//timeDisplay.Get("style").Set("fontFamily", "monospace")

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

}

func PlayerConstruct(window *wm.RiwoWindow) {
	themeKey := "blue"
	fg := wm.ThemeMap[themeKey]["vivid"]
	bg := wm.ThemeMap[themeKey]["faded"]
	mg := wm.ThemeMap[themeKey]["normal"]

	container := wm.Create().
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", bg)

	audio := wm.CreateKnown("audio")
	audio.DOM().Set("src", "about:blank")

	timeDisplay := wm.
		Create().
		Text("00:00 | 00:00").
		Style("color", fg).
		Style("textAlign", "center").
		Style("marginBottom", "10px")

	controls := wm.
		Create().
		Style("display", "flex").
		Style("alignItems", "center")

	controls2 := wm.
		Create().
		Style("display", "flex").
		Style("alignItems", "center")

	styleButton := func(b *wm.RiwoElement) {
		b.Style("padding", "5px 15px").
			Style("backgroundColor", bg).
			Style("color", "black").
			Style("border", "solid "+mg).
			Style("borderRadius", "0").
			Style("cursor", "url(assets/cursor-inverted.svg), auto").
			Style("font-family", "monospace").
			Style("font-weight", "bold")
		b.Callback("mouseover", func(this js.Value, args []js.Value) interface{} {
			b.Style("background", fg)
			b.Style("color", bg)
			return nil
		})
		b.Callback("mouseout", func(this js.Value, args []js.Value) interface{} {
			b.Style("background", bg)
			b.Style("color", "black")
			return nil
		})
	}

	prevButton := wm.Create().Text("[<<]")
	playButton := wm.Create().Text("[|>]")
	nextButton := wm.Create().Text("[>>]")

	volUpButton := wm.Create().Text("(+)")
	volDnButton := wm.Create().Text("(-)")

	fileButton := wm.Create().Text("FILE")

	styleButton(prevButton)
	styleButton(playButton)
	styleButton(nextButton)
	styleButton(volDnButton)
	styleButton(volUpButton)
	styleButton(fileButton)

	fileStatus := wm.Create().
		Style("color", mg).
		Style("margin", "10px").
		Text("No files loaded")

	fileInput := wm.
		CreateKnown("input").
		Style("display", "none")
	fileInput.DOM().Set("type", "file")
	fileInput.DOM().Set("accept", "*/*")

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

	prevButtonUp := func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.DOM().Get("currentTime").Float()
		audio.DOM().Set("currentTime", currentTime-10)
		return nil
	}
	nextButtonUp := func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.DOM().Get("currentTime").Float()
		audio.DOM().Set("currentTime", currentTime+10)
		return nil
	}
	prevButton.Callback("mouseup", prevButtonUp)
	nextButton.Callback("mouseup", nextButtonUp)

	playButtonUp := func(this js.Value, args []js.Value) interface{} {
		if audio.DOM().Get("paused").Bool() {
			audio.DOM().Call("play")
			playButton.Text("[||]")
		} else {
			audio.DOM().Call("pause")
			playButton.Text("[|>]")
		}
		return nil
	}

	volDnUp := func(this js.Value, args []js.Value) interface{} {
		currentVol := audio.DOM().Get("volume").Float()
		if currentVol > 0.1 {
			audio.DOM().Set("volume", currentVol-0.1)
		}
		return nil
	}
	volUpUp := func(this js.Value, args []js.Value) interface{} {
		currentVol := audio.DOM().Get("volume").Float()
		if currentVol < 1.0 {
			audio.DOM().Set("volume", currentVol+0.1)
		}
		return nil
	}

	fileUp := func(this js.Value, args []js.Value) interface{} {
		fileInput.DOM().Call("click")
		return nil
	}

	audioUpdate := func(this js.Value, args []js.Value) interface{} {
		currentTime := audio.DOM().Get("currentTime").Float()
		duration := audio.DOM().Get("duration").Float()
		if !js.Global().Call("isNaN", duration).Bool() {
			timeDisplay.DOM().Set("textContent", formatTime(currentTime)+" | "+formatTime(duration))
		}
		if currentTime == duration {
			audio.DOM().Call("pause")
			playButton.Text("[|>]")
		}
		return nil
	}
	fileInputChange := func(this js.Value, args []js.Value) interface{} {
		files := fileInput.DOM().Get("files")
		if files.Length() > 0 {
			file := files.Index(0)
			url := js.Global().Get("URL").Call("createObjectURL", file)
			audio.DOM().Set("src", url)
			fileStatus.Text(":: " + file.Get("name").String())
			playButton.Text("[|>]")
		}
		return nil
	}

	fileInput.Callback("change", fileInputChange)
	audio.Callback("timeupdate", audioUpdate)
	fileButton.Callback("mouseup", fileUp)
	volDnButton.Callback("mouseup", volDnUp)
	volUpButton.Callback("mouseup", volUpUp)
	playButton.Callback("mouseup", playButtonUp)
	prevButton.Callback("mouseup", prevButtonUp)
	nextButton.Callback("mouseup", nextButtonUp)

	controls.Append(prevButton, playButton, nextButton)
	controls2.Append(volDnButton, volUpButton, fileButton)

	container.Append(timeDisplay, controls, fileStatus, controls2, fileInput, audio)

	window.Content.
		Inner("").
		Append(container)
}
