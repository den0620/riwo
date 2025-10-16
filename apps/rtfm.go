package apps

import (
	//"fmt"
	"riwo/wm"
	//"strconv"
	//"syscall/js"
)

func init() {
	// Read This Fabulous Manual
	AppRegistry["RTFM"] = rtfmConstruct
	/*
	  I want to move basic info/manual
	from developer tools -> console to
	a preopened ~fullscreen "RTFM" app
	that would contain instructions on
	basic R(W)IO actions and maybe doc
	*/
}

func rtfmConstruct(window *wm.RiwoWindow) {
	themeKey := "aqua"
	//fg := wm.ThemeMap[themeKey]["normal"]
	mg := wm.ThemeMap[themeKey]["vivid"]
	bg := wm.ThemeMap[themeKey]["faded"]

	container := wm.Create()
	container.
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", bg)
	exampleTitle := wm.Create()
	exampleTitle.
		Text("blah blah todo").
		Style("fontSize", "1.5em").
		Style("color", mg)
	container.Append(exampleTitle)

	/*
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

	themesPanel := wm.Create()
	themesPanel.
		Style("display", "grid").
		Style("gridTemplateColumns", "repeat(4, 1fr)").
		Style("gap", "8px").
		Style("maxWidth", "400px").
		Style("margin", "0 auto")


	for key, theme := range wm.ThemeMap {
		if wm.Verbose {
			wm.Print(fmt.Sprintf("\tTheme[%s] -> %v", key, theme))
		}
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
			Listen("mouseover", over).
			Listen("mouseout", out).
			Listen("mousedown", click).
			Mount(themesPanel)
	}

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


	updateClock.Invoke()
	*/
	window.Content.Inner("").Append(container)
}

