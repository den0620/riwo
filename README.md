# riwo
small webassembly rio-like implementation in go

riwo stands for "Riwo is web one"

hardly inspired by Plan 9's [Rio](https://9p.io/wiki/plan9/using_rio/index.html)

## How it works

It uses [syscall/js](https://pkg.go.dev/syscall/js) to manipulate DOM and mimic rio

## Roadmap

- [x] New
- [x] Resize
- [x] Move
- [x] Delete
- [x] Hide

## Possible known issues

There may be some inefficient eventListeners

Menu opens with single RMB click and NOT hold because I found it simpler
