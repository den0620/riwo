package apps

import (
	//"fmt"
	"riwo/wm"
	//"strconv"
	//"syscall/js"
)

func init() {
	AppRegistry["Mahjongg"] = mahjonggConstruct
}

type mahjonggBrick struct {
	Type	int				// Brick type, identical are deleteable
	Content	*wm.RiwoObject	// Connected DOM element.
}
func createMahjonggBrick(brickType int, mgColor string) *mahjonggBrick {
	return &mahjonggBrick{
		Type:    brickType,
		Content: wm.Create().
			Text("🍔").
			Style("height", "3rem").
			Style("width", "2rem").
			Style("justifyContent", "center").
			Style("alignItems", "center").
			Style("backgroundColor", mgColor),
	}
}

func mahjonggConstruct(window *wm.RiwoWindow) {
	//fg := wm.ThemeMap["yellow"]["normal"]
	mg := wm.ThemeMap["yellow"]["faded"]
	bg := wm.ThemeMap["yellow"]["vivid"]

	container := wm.Create()
	container.
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", bg)
	
	examplebrick1 := createMahjonggBrick(1, mg)
	
	container.Append(examplebrick1.Content)

	// ...

	window.Content.Inner("").Append(container)
}

/*
func applyTheme(e *wm.RiwoObject, theme map[string]string) {
	e.
		Style("cursor", wm.CursorInvertUrl).
		Style("padding", "10px, 20px").
		Style("color", "#000000").
		Style("backgroundColor", theme["faded"]).
		Style("border", "solid "+theme["vivid"]).
		Style("borderRadius", 0)
}
*/
