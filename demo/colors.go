package demo


import (
  "syscall/js"
  "strconv"
)


var (
  verbose        bool = false
)


func Print(value string) {
  js.Global().Get("console").Call("log", value)
}


func Colors_StartInWindow(windowId int) {
  document := js.Global().Get("document")
  window := document.Call("getElementById", strconv.Itoa(windowId))
  window.Set("innerHTML", `<div style="width: 100%; height: 100%; margin: 0; padding: 0;">
  <div title="White" style="background-color: #ffffff; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Gray" style="background-color: #777777; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Red" style="background-color: #bb5d5d; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Green" style="background-color: #88cc88; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Blue" style="background-color: #00aaff; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Yellow" style="background-color: #eeee9e; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Light Blue" style="background-color: #9eeeee; width: 12.5%; height: 100%; float: left;"></div>
  <div title="Light Gray" style="background-color: #cccccc; width: 12.5%; height: 100%; float: left;"></div>
</div>`)
  return
}


