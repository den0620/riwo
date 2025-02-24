# THESE DOCS ARE OUT OF DATE

This was done primary for me to navigate in my shit

To be done

# How exactly does it work?

idk

# Window structure

Window is Go type with these entries:
```
type Window struct {
	ID      int      // For the most part unites DOM object and Go object
	Element js.Value // Connected DOM element.
	// Tho "Move", "Resize", "Delete" and "Hide" are basic ones
	ContextEntries []struct {
		name     string
		callback func()
	}
}
```

ID is window id (how sudden)

Element is what it looks like as a html object (not what inside, whole window)

So you only should utilize innerHTML

## Handle exiting

If ParticularWindow.ID gets `-1` then it was deleted and app should stop

## Assigning Context Menu to a window

```go
// Implies existance of func ParticularCallback(window *Window)
// And ParticularWindow of type `Window`

customEntries = map[string]js.Value

customEntries["MyEntry"] = ParticularCallback(ParticularWindow)

ParticularWindow.ContextEntries = customEntries

```

## Accessing Standard Colors

Standard colors (those that are in readme) are available in `riwo/wm` as `GetColor` nested map variable:

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


