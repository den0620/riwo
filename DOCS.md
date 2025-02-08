This was done primary for me to navigate in my shit

# How exactly does it work?

Look, there is a `main.go` file and it contains some initial scripts to put everything together and there is wm folder that contains all window manager parts that were split into separate files from original `wm/init.go` in post-`v1.5` (that was split from `main.go` in `v1.0` and `v1.5`) and they generate some abstact type `Window` that containes [all shit](#window-structure) and [app](#developing-apps) inside apps/ gets this object at its disposal that suggests \<div\> manipulation and custom logic and [custom context menu](#assigning-context-menu-to-a-window) but wm itself manages window actions (and basic ones) and background and context menu and all this loads with `assets/wasm_exec.js` by Go team

![Sorry i dont know whos the artist]()

# Window structure

Window is Go type with these entries:
```
type Window struct {
	ID      int
	Element js.Value // Connected DOM element.
	ContextEntries {} map // Tho "Move", "Resize", "Delete" and "Hide" are basic ones
}
```

ID is window id (how sudden)

Element is what it looks like as a html object (now what inside, whole window)

ContextEntries is map that should look like this:

## Assigning Context Menu to a window

```go
// Implies existance of func ParticularCallback(window *Window)
// And ParticularWindow of type `Window`

ParticularContextMenu = map[string]js.Value

ParticularContextMenu["MyEntry"] = ParticularCallback(ParticularWindow)

ParticularWindow.ContextEntries = ParticularWindowsContextMenu

```

## Accessing Standard Colors

Standard colors (those that are in readme) are available in `riwo/wm` as GetColor nested map variable:

```
color string = GetColor["aqua"]["vivid"]
```

will get you color "#8888cc" of type `string`

### All Colors:
- monochrome
- red
- green
- blue
- yellow
- aqua
- gray
### All Subcolors:
- faded
- normal
- vivid

# Developing apps

You should check out whole [Window](#window-structure) section and write single(or I dont know, maybe multiple)-file program that imports `Window` type and does its business and put it in `apps/`, then register should find it and make available from context menu


