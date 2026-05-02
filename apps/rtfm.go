package apps

import "riwo/wm"


func RTFMConstruct(window *wm.RiwoWindow) {
	window.Title = "Read This Fabulous Manual"

	mg := wm.GetBorderColorStr("aqua")
	bg := wm.GetBackgroundColorStr("aqua")
	faded := wm.GetTheme("aqua")

	scroll := wm.Create()
	scroll.
		Style("height", "100%").
		Style("overflowY", "auto").
		Style("overflowX", "hidden").
		Style("boxSizing", "border-box").
		Style("padding", "12px 16px").
		Style("backgroundColor", bg).
		Style("color", "#111111").
		Style("lineHeight", "1.45").
		Style("fontSize", "15px")

	h1 := wm.Create().
		Text("Riwo basics").
		Style("fontSize", "1.35em").
		Style("color", mg).
		Style("margin", "0 0 8px 0").
		Style("fontWeight", "bold")

	blurb := wm.Create().
		Text("Riwo is a Rio-style window shell in the browser: the window manager draws and launcher tiles; apps are separate WebAssembly binaries. This page lists the gestures and modes you will use.")

	prologue := wm.Create().
		Text("Pointers: LMB = left mouse button, RMB = right mouse button. Most actions assume you use RMB unless noted.")

	sections := [][]string{
		{
			"Mouse and global behavior",
			"LMB anywhere cancels whatever mode or menu is pending (closes the context menu, drops Move/Resize selection, clears New/Delete/Hide intentions).",
			"Click app's body so the window receives focus and stacks to the front.",
		},
		{
			"Context menu (RMB)",
			"Click (not hold) RMB on the empty workspace to open the context menu. Choose an operation with another RMB on the label you want.",
			"Rows include Move - hold RMB on a window and drag to move it; Resize - RMB-select a window, then RMB-drag a region from which the new rectangle is taken; New - drag a rectangle to spawn a blank window with apps to choose from; Delete - RMB-select a window to delete; Hide - tuck a window and list it later at the bottom of the context menu.",
		},
		{
			"Windows and launcher",
			"A new desktop window receives the launcher grid listing every guest WASM app. Pick a tile to load that binary into this frame only.",
			"When a guest is running, the context menu (RMB after focusing that window) includes Exit. It restores the launcher grid and clears kernel guest hooks.",
		},
		{
			"Per-app menus (RMB on content)",
			"Guests can register optional entries (Mahjong / ZClock, ...). Appear underneath the defaults when that window has focus.",
		},
		{
			"Debugging",
			"Open DevTools · Console. Use Logging() to toggle verbose wm logs exported on the JS global.",
		},
	}

	body := wm.Create()
	body.Append(h1, blurb, prologue)

	for _, sec := range sections {
		header := wm.Create().
			Text(sec[0]).
			Style("fontSize", "1.08em").
			Style("color", mg).
			Style("marginTop", "14px").
			Style("marginBottom", "6px").
			Style("fontWeight", "bold")

		ul := wm.CreateKnown("ul").Style("margin", "4px 0 8px 0").Style("paddingLeft", "1.25rem")
		for _, line := range sec[1:] {
			li := wm.CreateKnown("li")
			li.Text(line).
				Style("marginBottom", "6px")
			ul.Append(li)
		}
		body.Append(header, ul)
	}

	footer := wm.Create().
		Text("Tip: if something breaks, reload the browser tab.").
		Style("marginTop", "14px").
		Style("fontSize", "0.92em").
		Style("color", faded["vivid"])

	body.Append(footer)

	scroll.Append(body)
	window.Content.Inner("").Append(scroll)
}

