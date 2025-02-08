# riwo
small webassembly rio-like window manager in go

riwo stands for "Riwo is web one"

hardly inspired by Plan 9's [Rio](https://9p.io/wiki/plan9/using_rio/index.html)

![Preview](screenshot.webp)

## How it works

IT WON'T BUILD AND RUN. IT'S IN PROCESS OF A HUGE REWORK. I JUST WANTED TO SAVE MIDDLE CHANGES. USE ONE FROM RELEASE.

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
- [ ] Window-specific menu
### Accessibility
- [x] Touch adaptation

## Possible known issues

There may be some inefficient eventListeners (tho i removed per-window ones)

Menu opens with single RMB click and NOT hold because I found it simpler

SVG cursors may be junky

Testing needed

