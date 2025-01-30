package rc


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
func GoVerbose(this js.Value, args []js.Value) interface{} {
  verbose = !verbose
  Print("Verbose : " + strconv.FormatBool(verbose))
  return nil
}


func StartInWindow(windowId int) {
  document := js.Global().Get("document")
  window := document.Call("getElementById", strconv.Itoa(windowId))
  window.Set("innerHTML", "<h3>Now it's rc (kinda)</h3><p>this window's id is "+strconv.Itoa(windowId)+"</p>")
}
