package apps

import (
	"riwo/wm"
	"syscall/js"
)

func init() {
	AppRegistry["Monaco"] = monacoConstruct
}

func monacoConstruct(window *wm.RiwoWindow) {
	window.Title = "Monaco Editor (minimal)"

	// Bad idea...
	iframe := wm.CreateKnown("iframe")
	iframe.
		Attr("src", "apps/Monaco/index.html").
		Style("width", "100%").
		Style("height", "100%").
		Style("border", "none")

	iframe.Listen("load", func(this js.Value, args []js.Value) interface{} {
		contentWindow := iframe.DOM().Get("contentWindow")
		if contentWindow.Truthy() {
			contentWindow.Call("postMessage",
				js.ValueOf(map[string]interface{}{
					"type":     "init",
					"content":  "", // <-- string expected
					"language": "markdown",
				}), "*")
		}
		return nil
	})

	window.Content.
		Inner("").
		Append(iframe)
}

