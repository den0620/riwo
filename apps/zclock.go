package apps

import (
	"fmt"
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func init() {
	AppRegistry["ZClock"] = ClockConstruct
}

func AppZClock(window *wm.Window) {
	document := js.Global().Get("document")
	isSettings := false
	currentTheme := "aqua"

	// Create main container
	container := document.Call("createElement", "div")
	container.Get("style").Set("height", "100%")
	container.Get("style").Set("display", "flex")
	container.Get("style").Set("flexDirection", "column")
	container.Get("style").Set("justifyContent", "center")
	container.Get("style").Set("alignItems", "center")
	container.Get("style").Set("backgroundColor", wm.ThemeMap["aqua"]["faded"])

	// Create clock display
	clockDisplay := document.Call("createElement", "div")
	clockDisplay.Get("style").Set("fontSize", "4em")
	clockDisplay.Get("style").Set("color", wm.ThemeMap["aqua"]["vivid"])

	// Settings container with simplified styling
	settingsContainer := document.Call("createElement", "div")
	settingsContainer.Get("style").Set("display", "none")
	settingsContainer.Get("style").Set("padding", "20px")
	settingsContainer.Get("style").Set("textAlign", "center")

	// Settings title
	settingsTitle := document.Call("createElement", "div")
	settingsTitle.Set("textContent", "Clock Settings")
	settingsTitle.Get("style").Set("fontSize", "1.5em")
	settingsTitle.Get("style").Set("marginBottom", "15px")
	settingsTitle.Get("style").Set("color", wm.ThemeMap["aqua"]["vivid"])

	// UTC adjustment section
	utcSection := document.Call("createElement", "div")
	utcSection.Get("style").Set("marginBottom", "15px")

	utcLabel := document.Call("createElement", "div")
	utcLabel.Set("textContent", "UTC Offset")
	utcLabel.Get("style").Set("marginBottom", "5px")
	utcLabel.Get("style").Set("cursor", "url(assets/cursor.svg), auto")

	// Custom number input container
	utcInputContainer := document.Call("createElement", "div")
	utcInputContainer.Get("style").Set("display", "flex")
	utcInputContainer.Get("style").Set("alignItems", "center")
	utcInputContainer.Get("style").Set("justifyContent", "center")
	utcInputContainer.Get("style").Set("gap", "10px")

	utcValue := document.Call("createElement", "div")
	utcValue.Set("textContent", "5")
	utcValue.Get("style").Set("minWidth", "30px")
	utcValue.Get("style").Set("textAlign", "center")

	// Function to style UTC buttons according to theme
	styleUtcButton := func(btn js.Value, theme string) {
		btn.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
		btn.Get("style").Set("padding", "10px 20px")
		btn.Get("style").Set("backgroundColor", wm.ThemeMap[theme]["faded"])
		btn.Get("style").Set("color", "black")
		btn.Get("style").Set("border", "solid "+wm.ThemeMap[theme]["vivid"])
		btn.Get("style").Set("borderRadius", "0")
	}

	utcDecrease := document.Call("createElement", "button")
	utcDecrease.Set("textContent", "-")
	styleUtcButton(utcDecrease, "aqua")

	utcIncrease := document.Call("createElement", "button")
	utcIncrease.Set("textContent", "+")
	styleUtcButton(utcIncrease, "aqua")

	// Theme section
	themeSection := document.Call("createElement", "div")

	themeLabel := document.Call("createElement", "div")
	themeLabel.Set("textContent", "Color Theme")
	themeLabel.Get("style").Set("marginBottom", "10px")
	themeLabel.Get("style").Set("cursor", "url(assets/cursor.svg), auto")

	// Theme buttons container
	themeContainer := document.Call("createElement", "div")
	themeContainer.Get("style").Set("display", "grid")
	themeContainer.Get("style").Set("gridTemplateColumns", "repeat(4, 1fr)")
	themeContainer.Get("style").Set("gap", "8px")
	themeContainer.Get("style").Set("maxWidth", "400px")
	themeContainer.Get("style").Set("margin", "0 auto")

	themes := []string{"monochrome", "red", "green", "blue", "yellow", "aqua", "gray"}

	// Create theme buttons
	for _, theme := range themes {
		themeBtn := document.Call("createElement", "div")
		themeBtn.Get("style").Set("padding", "10px")
		themeBtn.Get("style").Set("textAlign", "center")
		themeBtn.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
		themeBtn.Get("style").Set("backgroundColor", wm.ThemeMap[theme]["faded"])
		themeBtn.Get("style").Set("color", "black")
		themeBtn.Get("style").Set("border", "solid "+wm.ThemeMap[theme]["vivid"])
		themeBtn.Get("style").Set("borderRadius", "0")
		themeBtn.Set("textContent", theme)

		// Theme selection handler
		themeBtn.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			themeBtn.Get("style").Set("backgroundColor", wm.ThemeMap[theme]["normal"])
			return nil
		}))
		themeBtn.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			themeBtn.Get("style").Set("backgroundColor", wm.ThemeMap[theme]["faded"])
			return nil
		}))
		themeBtn.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			newTheme := this.Get("textContent").String()
			currentTheme = newTheme
			container.Get("style").Set("backgroundColor", wm.ThemeMap[newTheme]["faded"])
			clockDisplay.Get("style").Set("color", wm.ThemeMap[newTheme]["vivid"])
			settingsTitle.Get("style").Set("color", wm.ThemeMap[newTheme]["vivid"])
			// Update UTC buttons style
			styleUtcButton(utcDecrease, newTheme)
			styleUtcButton(utcIncrease, newTheme)
			return nil
		}))

		themeContainer.Call("appendChild", themeBtn)
	}

	// UTC button handlers
	utcDecrease.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		utcDecrease.Get("style").Set("backgroundColor", wm.ThemeMap[currentTheme]["normal"])
		return nil
	}))
	utcDecrease.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		utcDecrease.Get("style").Set("backgroundColor", wm.ThemeMap[currentTheme]["faded"])
		return nil
	}))
	utcDecrease.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		current, _ := strconv.Atoi(utcValue.Get("textContent").String())
		if current > -12 {
			utcValue.Set("textContent", strconv.Itoa(current-1))
		}
		return nil
	}))

	utcIncrease.Call("addEventListener", "mouseover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		utcIncrease.Get("style").Set("backgroundColor", wm.ThemeMap[currentTheme]["normal"])
		return nil
	}))
	utcIncrease.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		utcIncrease.Get("style").Set("backgroundColor", wm.ThemeMap[currentTheme]["faded"])
		return nil
	}))
	utcIncrease.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		current, _ := strconv.Atoi(utcValue.Get("textContent").String())
		if current < 14 {
			utcValue.Set("textContent", strconv.Itoa(current+1))
		}
		return nil
	}))

	// Append UTC elements
	utcInputContainer.Call("appendChild", utcDecrease)
	utcInputContainer.Call("appendChild", utcValue)
	utcInputContainer.Call("appendChild", utcIncrease)

	utcSection.Call("appendChild", utcLabel)
	utcSection.Call("appendChild", utcInputContainer)

	// Append theme elements
	themeSection.Call("appendChild", themeLabel)
	themeSection.Call("appendChild", themeContainer)

	// Append all settings elements
	settingsContainer.Call("appendChild", settingsTitle)
	settingsContainer.Call("appendChild", utcSection)
	settingsContainer.Call("appendChild", themeSection)

	// Helper function to format time components
	formatTime := func(n int) string {
		if n < 10 {
			return "0" + strconv.Itoa(n)
		}
		return strconv.Itoa(n)
	}

	// Update function
	var updateClock js.Func
	updateClock = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := js.Global().Get("Date").New()
		utcOffset, _ := strconv.Atoi(utcValue.Get("textContent").String())

		hours := (now.Call("getUTCHours").Int() + utcOffset + 24) % 24
		minutes := now.Call("getUTCMinutes").Int()
		seconds := now.Call("getUTCSeconds").Int()

		timeStr := formatTime(hours) + ":" + formatTime(minutes) + ":" + formatTime(seconds)
		clockDisplay.Set("textContent", timeStr)

		// Schedule next update
		js.Global().Call("setTimeout", updateClock, 1000)
		return nil
	})

	// Add settings toggle to window context menu
	window.ContextEntries = []wm.ContextEntry{
		{
			Name: "Settings",
			Callback: func() {
				isSettings = !isSettings
				if isSettings {
					settingsContainer.Get("style").Set("display", "block")
					clockDisplay.Get("style").Set("display", "none")
				} else {
					settingsContainer.Get("style").Set("display", "none")
					clockDisplay.Get("style").Set("display", "block")
				}
				if wm.Verbose {
					wm.Print("zclock settings toggled: " + strconv.FormatBool(isSettings))
				}
			},
		},
	}

	// Append all elements
	container.Call("appendChild", clockDisplay)
	container.Call("appendChild", settingsContainer)

	// Clear window and add container
	window.DOM.Set("innerHTML", "")
	window.DOM.Call("appendChild", container)

	// Start the clock
	updateClock.Invoke()
}

func ClockConstruct(window *wm.Window) {
	fg := wm.ThemeMap["aqua"]["normal"]
	mg := wm.ThemeMap["aqua"]["vivid"]
	bg := wm.ThemeMap["aqua"][""]

	container := wm.Create()
	container.
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", bg)

	clock := wm.Create()
	clock.
		Style("fontSize", "4em").
		Style("color", fg)

	// simple container styling :D
	settings := wm.Create()
	settings.
		Style("display", "none").
		Style("padding", "20px").
		Style("textAlign", "center")

	settingsTitle := wm.Create()
	settingsTitle.
		Text("Clock Settings").
		Style("fontSize", "1.5em").
		Style("marginBottom", "15px").
		Style("color", mg)

	settingsUtc := wm.Create()
	settingsUtc.Style("marginBottom", "15px")

	utc := wm.Create().
		Style("display", "flex").
		Style("alignItems", "center").
		Style("justifyContent", "center").
		Style("gap", "10px")

	utcLabel := wm.Create()
	utcLabel.
		Text("UTC Offset").
		Style("cursor", "url(assets/cursor.svg), auto").
		Style("marginBottom", "5px")

	utcInput := wm.Create()
	utcInput.
		Text("7").
		Style("minWidth", "30px").
		Style("textAlign", "center")

	utcHourIncrase := wm.Create()
	utcHourDecrase := wm.Create()

	utcHourIncrase.
		Style("cursor", "url(assets/cursor-inverted.svg), auto").
		Style("padding", "10px, 20px").
		Style("color", "#000000").
		Style("backgroundColor", bg).
		Style("border", "solid "+mg).
		Style("borderRadius", 0).
		Text("+")
	utcHourDecrase.
		Style("cursor", "url(assets/cursor-inverted.svg), auto").
		Style("padding", "10px, 20px").
		Style("color", "#000000").
		Style("backgroundColor", bg).
		Style("border", "solid "+mg).
		Style("borderRadius", 0).
		Text("-")

	themesPanel := wm.Create()
	themesPanel.
		Style("display", "grid").
		Style("gridTemplateColumns", "repeat(4, 1fr)").
		Style("gap", "8px").
		Style("maxWidth", "400px").
		Style("margin", "0 auto")

	themeKey := "aqua"

	for key, theme := range wm.ThemeMap { // <-- not sure
		wm.Print(fmt.Sprintf("\tTheme[%s] -> %v", key, theme))
		themeButton := wm.Create()

		applyTheme(themeButton, theme)

		out := func(this js.Value, args []js.Value) interface{} {
			themeButton.Style("backgroundColor", wm.ThemeMap[key]["faded"])
			return nil
		}
		over := func(this js.Value, args []js.Value) interface{} {
			themeButton.Style("backgroundColor", wm.ThemeMap[key]["normal"])
			return nil
		}

		click := func(this js.Value, args []js.Value) interface{} {
			newTheme := this.Get("textContent").String()

			container.Style("backgroundColor", wm.ThemeMap[newTheme]["faded"])
			clock.Style("color", wm.ThemeMap[newTheme]["vivid"])
			settingsTitle.Style("color", wm.ThemeMap[newTheme]["vivid"])

			applyTheme(utcHourDecrase, wm.ThemeMap[newTheme])
			applyTheme(utcHourIncrase, wm.ThemeMap[newTheme])

			themeKey = newTheme
			return nil
		}
		themeButton.
			Text(key).
			Callback("mouseover", over).
			Callback("mouseout", out).
			Callback("mousedown", click).
			Mount(themesPanel)
	}
	// Poop the bank
	decrDecorateMouseOver := func(this js.Value, args []js.Value) interface{} {
		utcHourDecrase.Style("backgroundColor", wm.ThemeMap[themeKey]["normal"])
		return nil
	}
	decrDecorateMouseOut := func(this js.Value, args []js.Value) interface{} {
		utcHourDecrase.Style("backgroundColor", wm.ThemeMap[themeKey]["faded"])
		return nil
	}
	decrClick := func(this js.Value, args []js.Value) interface{} {
		current, _ := strconv.Atoi(utcInput.DOM().Get("textContent").String())
		if current > -12 {
			utcInput.Text(strconv.Itoa(current - 1))
		}
		return nil
	}
	incrDecorateMouseOver := func(this js.Value, args []js.Value) interface{} {
		utcHourIncrase.Style("backgroundColor", wm.ThemeMap[themeKey]["normal"])
		return nil
	}
	incrDecorateMouseOut := func(this js.Value, args []js.Value) interface{} {
		utcHourIncrase.Style("backgroundColor", wm.ThemeMap[themeKey]["faded"])
		return nil
	}
	incrClick := func(this js.Value, args []js.Value) interface{} {
		current, _ := strconv.Atoi(utcInput.DOM().Get("textContent").String())
		if current > -12 {
			utcInput.Text(strconv.Itoa(current + 1))
		}
		return nil
	}

	// i'M rEaLLy NoT OkaY

	utcHourDecrase.
		Callback("mouseover", decrDecorateMouseOver).
		Callback("mouseout", decrDecorateMouseOut).
		Callback("mousedown", decrClick)
	utcHourIncrase.
		Callback("mouseover", incrDecorateMouseOver).
		Callback("mouseout", incrDecorateMouseOut).
		Callback("mousedown", incrClick)

	// not 9:30 -> 09:30 is better
	formatTime := func(n int) string {
		if n < 10 {
			return "0" + strconv.Itoa(n)
		}
		return strconv.Itoa(n)
	}

	var updateClock js.Func
	updateClock = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := js.Global().Get("Date").New()
		utcOffset, _ := strconv.Atoi(utcInput.DOM().Get("textContent").String())

		hours := (now.Call("getUTCHours").Int() + utcOffset + 24) % 24
		minutes := now.Call("getUTCMinutes").Int()
		seconds := now.Call("getUTCSeconds").Int()

		timeStr := formatTime(hours) + ":" + formatTime(minutes) + ":" + formatTime(seconds)
		clock.Text(timeStr)

		// Schedule next update
		js.Global().Call("setTimeout", updateClock, 1000)
		return nil
	})

	isSettingsShown := false

	window.ContextEntries = []wm.ContextEntry{
		{
			Name: "Settings",
			Callback: func() {
				isSettingsShown = !isSettingsShown
				if isSettingsShown {
					settings.Style("display", "block")
					clock.Style("display", "none")
				} else {
					settings.Style("display", "none")
					clock.Style("display", "block")
				}
				if wm.Verbose {
					wm.Print("zclock settings toggled: " + strconv.FormatBool(isSettingsShown))
				}
			},
		},
	}
	utc.Append(utcHourDecrase, utcInput, utcHourIncrase)
	settingsUtc.Append(utcLabel, utc)

	settings.Append(settingsUtc, themesPanel)
	container.Append(clock, settings)
	window.DOM.Set("innerHTML", "")
	window.DOM.Call("appendChild", container.DOM())

	updateClock.Invoke()
}

func applyTheme(e *wm.RiwoElement, theme map[string]string) {
	e.
		Style("cursor", "url(assets/cursor-inverted.svg), auto").
		Style("padding", "10px, 20px").
		Style("color", "#000000").
		Style("backgroundColor", theme["vivid"]).
		Style("border", "solid "+theme["vivid"]).
		Style("borderRadius", 0)
}
