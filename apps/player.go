package apps

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
)

func init() {
	AppRegistry["DPlayer"] = PlayerConstruct
}

func PlayerConstruct(window *wm.RiwoWindow) {
	fg := wm.GetBorderColorStr("blue")
	bg := wm.GetBackgroundColorStr("blue")
	mg := wm.GetFontColorStr("blue")

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

	styleButton := func(b *wm.RiwoObject) {
		b.Style("padding", "5px 15px").
			Style("backgroundColor", bg).
			Style("color", "black").
			Style("border", "solid "+mg).
			Style("borderRadius", "0").
			Style("cursor", wm.CursorInvertUrl).
			Style("font-family", "monospace").
			Style("font-weight", "bold")
		b.Listen("mouseover", func(this js.Value, args []js.Value) interface{} {
			b.Style("background", fg)
			b.Style("color", bg)
			return nil
		})
		b.Listen("mouseout", func(this js.Value, args []js.Value) interface{} {
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
	prevButton.Listen("mouseup", prevButtonUp)
	nextButton.Listen("mouseup", nextButtonUp)

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

	fileInput.Listen("change", fileInputChange)
	audio.Listen("timeupdate", audioUpdate)
	fileButton.Listen("mouseup", fileUp)
	volDnButton.Listen("mouseup", volDnUp)
	volUpButton.Listen("mouseup", volUpUp)
	playButton.Listen("mouseup", playButtonUp)
	prevButton.Listen("mouseup", prevButtonUp)
	nextButton.Listen("mouseup", nextButtonUp)

	controls.Append(prevButton, playButton, nextButton)
	controls2.Append(volDnButton, volUpButton, fileButton)

	container.Append(timeDisplay, controls, fileStatus, controls2, fileInput, audio)

	window.Content.
		Inner("").
		Append(container)
}
