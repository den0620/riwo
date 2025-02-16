# riwo
Small webassembly rio-like window manager in go

riwo stands for "Riwo is web one"

hardly inspired by Plan 9's [Rio](https://9p.io/wiki/plan9/using_rio/index.html)

![Preview](screenshot.webp)

## How it works

It uses [syscall/js](https://pkg.go.dev/syscall/js) to manipulate DOM and mimic rio

Windows are \<div\>s with html inside

Try it here: [deployment](https://ninefid.uk.to/riwo)

## Roadmap
### Actions
- [x] New
- [x] Resize
- [x] Move
- [x] Delete
- [x] Hide
- [x] Window-specific menu
- [x] Bearable apps
### Accessibility
- [ ] Touch adaptation (broke for some reason)

## Possible known issues

Menu opens with single RMB click and NOT hold because I found it simpler

SVG cursors may be junky

APP_defaults' buttons have default cursors (will fix soon)

Hidden windows appear in context menu unsorted (will fix soon)

Context menu appears after modes (will fix soon)

If mode was interrupted without mouseup things may brake (I really dont want to fix this)

