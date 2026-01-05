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

I would like to see any contribution

## Roadmap
### Actions
- [x] New
- [x] Resize
- [x] Move
- [x] Delete
- [x] Hide
- [x] Window-specific context menu entries
- [x] Bearable apps
- [ ] Make apps load as modules, not monolithic wasm
### Accessibility
- [x] Touch adaptation
### Default apps
- [x] Starter (`Default`)
- [x] Clock (`ZClock`)
- [x] Audio player (`DPlayer`)
- [x] Mahjongg (`Mahjongg`)
- [x] Monaco Editor (`Monaco`)
- [ ] Manual (`RTFM`)
- [ ] Gallery (`?`)
- [ ] Drawterm (`?`)
- [ ] Doom (`Doom`?)
- [ ] BoxedWine (`?`)
- [ ] Deus Ex Demo (`DXdemo`?)

## Possible known issues

Menu opens with single RMB click and NOT hold because I found it simpler

Buttons can be clicked with both RMB and LMB and NOT mousewheel button because I found it simpler

SVG cursors may be junky

If mode was interrupted without mouseup things may brake (I really dont want to fix this)

Mahjongg has no plan9-ish cursor for cursor "not allowed"

Mahjongg may not fit in phone's screen

<a href="https://star-history.com/#den0620/riwo&Date">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=den0620/riwo&type=Date&theme=dark" />
        <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=den0620/riwo&type=Date" />
        <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=den0620/riwo&type=Date" />
    </picture>
</a>

## Building

```shell
# Default build (outputs to build/main.wasm)
make all

# Custom output path
make all OUTPUT=build/example_dir/main.wasm

# Remove build
make clean
```

