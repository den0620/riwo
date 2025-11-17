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

	// ***
	
	window.Content.Inner("").Append(container)
}

