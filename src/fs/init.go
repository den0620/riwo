package fs


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


func InitializeStructure() {
  return
}
