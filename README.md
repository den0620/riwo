# riwo
Small webassembly rio-like window manager in go

riwo stands for "Riwo is web one"

hardly inspired by Plan 9's [Rio](https://9p.io/wiki/plan9/using_rio/index.html)

![Preview](screenshot.webp)

## How it works

It uses [syscall/js](https://pkg.go.dev/syscall/js) to manipulate DOM and mimic rio

Windows are \<div\>s with html inside

Better documentation soon™️

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
- [x] Touch adaptation (broke for some reason)
### Default apps
- [x] Starter (`Default`)
- [x] Clock (`ZClock`)
- [x] Audio player (`DPlayer`)
- [ ] Gallery ()
- [ ] Themer ()
- [ ] Mahjongg ()

## Possible known issues

Menu opens with single RMB click and NOT hold because I found it simpler

SVG cursors may be junky

If mode was interrupted without mouseup things may brake (I really dont want to fix this)

Apps or their processes may remain alive (this should be fixed by checking if underlying window isn't nil but idk)


<a href="https://star-history.com/#den0620/riwo&Date">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=den0620/riwo&type=Date&theme=dark" />
        <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=den0620/riwo&type=Date" />
        <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=den0620/riwo&type=Date" />
    </picture>
</a>
