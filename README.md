# riwo
Small webassembly rio-like window manager in go

riwo stands for "Riwo is web one"

hardly inspired by Plan 9's [Rio](https://9p.io/wiki/plan9/using_rio/index.html)

![Preview](screenshot.webp)

## How it works

It uses [syscall/js](https://pkg.go.dev/syscall/js) to manipulate DOM and mimic rio

Windows are \<div\>s with html inside

Credits to Go team for their [Go fonts](https://go.dev/blog/go-fonts)

Try it here: [deployment](https://ninefid.uk.to/riwo)

## Roadmap
### Actions
- [x] New
- [x] Resize
- [x] Move
- [x] Delete
- [x] Hide
- [x] Window-specific context menu entries
- [x] Bearable apps
### Accessibility
- [ ] Touch adaptation (broke for some reason)

## Possible known issues

Menu opens with single RMB click and NOT hold because I found it simpler

SVG cursors may be junky

If mode was interrupted without mouseup things may brake (I really dont want to fix this)

Apps or their processes may remain alive (this should be fixed by checking if underlying window isn't nil but idk)
