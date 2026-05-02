package wm

import "syscall/js"

// GuestMenuEntry is a rio-style context row registered with the JS kernel so WM wasm can show it while the owning app runs separately
type GuestMenuEntry struct {
	Name     string
	Callback func()
}

// RegisterGuestContextMenus publishes menu rows to the JS kernel (__riwoGuestBootstrap host)
func RegisterGuestContextMenus(windowID int, entries []GuestMenuEntry) {
	k := js.Global().Get("__riwoKernel")
	if !k.Truthy() {
		return
	}
	for _, ent := range entries {
		name := ent.Name
		cb := ent.Callback
		h := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if cb != nil {
				cb()
			}
			return nil
		})
		k.Call("guestContextMenuAppend", windowID, name, h)
	}
}

func KernelGuestContextMenuTitles(windowID int) []string {
	k := js.Global().Get("__riwoKernel")
	if !k.Truthy() {
		return nil
	}
	arr := k.Call("guestContextMenuTitles", windowID)
	if arr.Type() == js.TypeNull || arr.Type() == js.TypeUndefined {
		return nil
	}
	n := arr.Get("length").Int()
	if n <= 0 {
		return nil
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, arr.Index(i).String())
	}
	return out
}

func KernelInvokeGuestContextMenu(windowID int, index int) {
	k := js.Global().Get("__riwoKernel")
	if !k.Truthy() {
		return
	}
	k.Call("guestContextMenuInvoke", windowID, index)
}
