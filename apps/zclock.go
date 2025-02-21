package apps

import (
    "riwo/wm"
    "syscall/js"
    "strconv"
)

func init() {
    AppRegistry["ZClock"] = APP_zclock
}

var (
    isSettings bool
)

func APP_zclock(window *wm.Window) {
    document := js.Global().Get("document")
    isSettings = false
    
    // Create main container
    container := document.Call("createElement", "div")
    container.Get("style").Set("height", "100%")
    container.Get("style").Set("display", "flex")
    container.Get("style").Set("flexDirection", "column")
    container.Get("style").Set("justifyContent", "center")
    container.Get("style").Set("alignItems", "center")
    container.Get("style").Set("backgroundColor", wm.GetColor["aqua"]["faded"])
    
    // Create clock display
    clockDisplay := document.Call("createElement", "div")
    clockDisplay.Get("style").Set("fontSize", "4em")
    clockDisplay.Get("style").Set("fontFamily", "monospace")
    clockDisplay.Get("style").Set("color", wm.GetColor["aqua"]["vivid"])
    
    // Settings container
    settingsContainer := document.Call("createElement", "div")
    settingsContainer.Get("style").Set("display", "none")
    
    // UTC adjustment input
    utcLabel := document.Call("createElement", "label")
	utcLabel.Get("style").Set("cursor", "url(assets/cursor.svg), auto")
    utcLabel.Set("textContent", "UTC Offset: ")
    utcInput := document.Call("createElement", "input")
	utcInput.Get("style").Set("cursor", "url(assets/cursor-selection.svg) 12 12, auto")
    utcInput.Set("type", "number")
    utcInput.Set("value", "5")
    utcInput.Set("min", "-12")
    utcInput.Set("max", "14")
    
    // Color theme selector
    themeLabel := document.Call("createElement", "label")
	themeLabel.Get("style").Set("cursor", "url(assets/cursor.svg), auto")
    themeLabel.Set("textContent", "Color Theme: ")
    themeSelect := document.Call("createElement", "select")
	themeSelect.Get("style").Set("cursor", "url(assets/cursor-inverted.svg), auto")
    
    themes := []string{"monochrome", "red", "green", "blue", "yellow", "aqua", "gray"}
    for _, theme := range themes {
        option := document.Call("createElement", "option")
        option.Set("value", theme)
        option.Set("textContent", theme)
        themeSelect.Call("appendChild", option)
    }
    themeSelect.Set("value", "aqua")
    
    // Append settings elements
    settingsContainer.Call("appendChild", utcLabel)
    settingsContainer.Call("appendChild", utcInput)
    settingsContainer.Call("appendChild", document.Call("createElement", "br"))
    settingsContainer.Call("appendChild", themeLabel)
    settingsContainer.Call("appendChild", themeSelect)
    
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
        utcOffset, _ := strconv.Atoi(utcInput.Get("value").String())
        
        hours := (now.Call("getUTCHours").Int() + utcOffset + 24) % 24
        minutes := now.Call("getUTCMinutes").Int()
        seconds := now.Call("getUTCSeconds").Int()
        
        timeStr := formatTime(hours) + ":" + formatTime(minutes) + ":" + formatTime(seconds)
        clockDisplay.Set("textContent", timeStr)
        
        // Schedule next update
        js.Global().Call("setTimeout", updateClock, 1000)
        return nil
    })
    
    // Theme change handler
    themeSelect.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        theme := themeSelect.Get("value").String()
        container.Get("style").Set("backgroundColor", wm.GetColor[theme]["faded"])
        clockDisplay.Get("style").Set("color", wm.GetColor[theme]["vivid"])
        return nil
    }))
    
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
    window.Element.Set("innerHTML", "")
    window.Element.Call("appendChild", container)
    
    // Start the clock
    updateClock.Invoke()
}
