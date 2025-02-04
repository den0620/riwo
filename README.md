# riwo
small webassembly rio-like window manager in go

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
- [ ] Window-specific menu
### Accessibility
- [x] Touch adaptation

## Possible known issues

There may be some inefficient eventListeners (tho i removed per-window ones)

Menu opens with single RMB click and NOT hold because I found it simpler

SVG cursors may be junky

Testing needed

## Color palette

If you want to make something in Plan 9 (9front) style, here is its common colors palette:

- Window: #ffffff
- Border: #9eeeee
- Active Border: #55aaaa
- Background: #777777
### Red
- Faded: #ffeaea
- Normal: #df9595
- Vivid: #bb5d5d
### Green
- Faded: #eaffea
- Normal: #88cc88
- Vivid: #448844
### Blue
- Faded: #c0eaff
- Normal: #00aaff
- Vivid: #0088cc
### Yellow
- Faded: #ffffea
- Normal: #eeee9e
- Vivid: #99994c
### Light Blue
- Faded: #eaffff
- Normal: #9eeeee
- Vivid: #8888cc
### Light Gray
- Faded: #eeeeee
- Normal: #cccccc
- Vivid: #888888

