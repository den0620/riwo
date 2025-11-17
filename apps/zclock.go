package apps

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func init() {
	AppRegistry["ZClock"] = clockConstruct
}

func clockConstruct(window *wm.RiwoWindow) {
	//fg := wm.ThemeMap["aqua"]["normal"] // not used so compiler will shit its pants
	mg := wm.ThemeMap["aqua"]["vivid"]
	bg := wm.ThemeMap["aqua"]["faded"]

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
		Style("color", mg)

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
		Style("cursor", wm.CursorDefaultUrl).
		Style("marginBottom", "5px")

	utcInput := wm.Create()
	utcInput.
	    Text(strconv.Itoa(-js.Global().Get("Date").New().Call("getTimezoneOffset").Int()/60)).
		Style("minWidth", "30px").
		Style("textAlign", "center")

	utcHourIncrase := wm.Create()
	utcHourDecrase := wm.Create()

	utcHourIncrase.
		Style("cursor", wm.CursorInvertUrl).
		Style("padding", "10px, 20px").
		Style("color", "#000000").
		Style("backgroundColor", bg).
		Style("border", "solid "+mg).
		Style("borderRadius", 0).
		Text("+")
	utcHourDecrase.
		Style("cursor", wm.CursorInvertUrl).
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
		wm.Print("\tTheme["+key+"] -> "+themeKey) // was "theme" but gave an error
		themeButton := wm.Create()

		applyThemeToButton(themeButton, theme)

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

			applyThemeToButton(utcHourDecrase, wm.ThemeMap[newTheme])
			applyThemeToButton(utcHourIncrase, wm.ThemeMap[newTheme])

			themeKey = newTheme
			return nil
		}
		themeButton.
			Text(key).
			Listen("mouseover", over).
			Listen("mouseout", out).
			Listen("mousedown", click).
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
		Listen("mouseover", decrDecorateMouseOver).
		Listen("mouseout", decrDecorateMouseOut).
		Listen("mousedown", decrClick)
	utcHourIncrase.
		Listen("mouseover", incrDecorateMouseOver).
		Listen("mouseout", incrDecorateMouseOut).
		Listen("mousedown", incrClick)

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
		utcOffset, _ := strconv.Atoi(utcInput.From("textContent").String())

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

	window.MenuEntries = []wm.ContextEntry{
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
				wm.Print("zclock settings toggled: " + strconv.FormatBool(isSettingsShown))
			},
		},
	}
	utc.Append(utcHourDecrase, utcInput, utcHourIncrase)
	settingsUtc.Append(utcLabel, utc)

	settings.Append(settingsUtc, themesPanel)
	container.Append(clock, settings)

	window.Content.Inner("").Append(container)

	updateClock.Invoke()
}

