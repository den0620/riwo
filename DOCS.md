# Riwo

Riwo is a Rio-inspired window shell in the browser: the **window manager** compiles to **WebAssembly**, and **applications** are separate guests that may be **Go wasm**, **`js-mount` modules** (no iframe—they render into the wm content pane), or **raw wasm**. A small **JavaScript kernel** owns load paths, guest instantiation, launcher metadata, and cross-instance bookkeeping (guest context menus bridged across wasm instances).

This document reflects the architecture after splitting from the earlier single `main.wasm` tree and introducing **guestgen** (`tools/guestgen`).

---

## Repository layout

| Path | Role |
|------|------|
| **`entry/wm`** | WM `main`: registers globals (`Logging`, `LaunchDefault`), starts `wm.InitializeContextMenu`, `InitializeGlobalMouseEvents`. |
| **`entry/guest/<name>/`** | One directory per launcher tile. May include **`main.go`** (Go guest), **`guest.manifest.json`** (display name, runtime, paths), or both depending on runtime. |
| **`tools/guestgen`** | Walks **`entry/guest/*`**, writes **`assets/generated-manifest.js`** (launcher list for the kernel) and **`build/go-guest-dirs.txt`** (which dirs get **`go build` wasm`). |
| **`wm/`** | Window manager DOM helpers, gestures, launcher, syscall bridges to **`__riwoKernel`**. |
| **`apps/`** | Guest UI constructors (`*_Construct`), imported only by Go **`entry/guest/<name>`** packages. |
| **`assets/kernel.js`** | **`__riwoKernel`**, WM start, serialized guest spawn, runtime-specific loaders (`go` / **`js-mount`** / **`wasm-raw`**), menu bridge. |
| **`assets/js-mounts/*.js`** | ESM modules with **`export async function mount(host, ctx)`** (or default function) for **`js-mount`** guests. |
| **`build/`** | **`wasm_exec.js`**, **`wm.wasm`**, **`<guest>.wasm`** for Go guests, **`go-guest-dirs.txt`**, **`generated-manifest.js`** is under **`assets/`** (not **`build/`**). |
| **`go/go_wasm_exec.js`** | Copied into **`build/wasm_exec.js`** for canonical Go wasm. |
| **`apps/Monaco/`** | Static Monaco assets (AMD loader + `vs/*`), loaded by **`assets/js-mounts/monaco-mount.js`**. |

Module path: **`module riwo`** (`go.mod`).

**Pure Go helpers** (no `syscall/js`, safe to unit test on the host):

| Path | Role |
|------|------|
| **`audiometa/`** | **`ParseTags([]byte)`** for DPlayer: **FLAC** (Vorbis comment), **MP4/M4A** (AAC or ALAC in MP4 — walks **`moov` → … → `ilst`**, reads **`©nam`**, **`©ART`**, **`©alb`**, falls back **`aART`** for artist), **MP3** (**ID3** + **ID3v1**). |

---

## JavaScript: `__riwoKernel` API

The global **`globalThis.__riwoKernel`** is defined in **`assets/kernel.js`**. The WM and Go guests reach it via **`syscall/js`** (`js.Global().Get("__riwoKernel")`).

| Method | Called from | Behavior |
|--------|-------------|----------|
| **`startWM('build/wm.wasm')`** | kernel (internal) | Instantiates WM wasm with the Go `importObject` and **`go.run`**. |
| **`spawnGuestApp(launchName, windowId, contentElement)`** | WM (`KernelSpawnGuest`) | Looks up **`launchName`** in **`__RIWO_GENERATED_MANIFEST`**, then loads **`go`**, **`js-mount`**, or **`wasm-raw`** as configured. Spawn work is serialized on an internal queue. |
| **`disposeGuestForWindow(windowId)`** | WM when a window is torn down or exiting a guest | Clears guest bookkeeping and guest context menu rows for that id. |
| **`listGuestApps()`** | WM launcher | Returns sorted **`launchName`** strings from the manifest. |
| **`guestContextMenuAppend(windowId, name, callback)`** | Go guests | Registers a menu row; **`callback`** is usually a **`syscall/js.Func`**. |
| **`guestContextMenuTitles(windowId)`** / **`guestContextMenuInvoke(windowId, index)`** | WM menu builder | Reads titles and invokes the chosen callback. |
| **`consumeGuestBootstrap()`** | Go guests (`wm.RunGuestApp`) | Returns **`{ windowId, pane }`** once for the pending Go spawn, or **`null`**. Replaces older single-slot globals so concurrent launches do not clash. |

**`js-mount`** modules receive **`ctx.kernel`** with **`resolveURL`**, **`disposeGuestForWindow`** only (extend here for shared JS affordances).

---

## Development workflow

- **Regenerate the guest manifest** whenever you add or rename under **`entry/guest/`**: run **`make default`**, **`make wm-wasm`**, or any target that refreshes **`$(GUESTGEN_STAMP)`**. The stamp depends on **`tools/guestgen/main.go`** and **`entry/guest/*/main.go`**, **`entry/guest/*/guest.manifest.json`**.
- **Serve over HTTP** (e.g. `python3 -m http.server` from the repo root, or any static server). Loading **`index.html`** via **`file://`** often breaks **`fetch()`** of wasm and dynamic **`import()`** for **`js-mount`** guests; treat a local HTTP origin as the baseline.
- **Unit tests** that must not import **`syscall/js`** belong in packages such as **`audiometa/`**. Packages under **`apps/`** and **`wm/`** are built for **`GOOS=js`** in normal development; use **`GOOS=js GOARCH=wasm go build ./entry/...`** to compile guests after changes.

---

## DPlayer (file guest) and tags

The music player (**`apps/player.go`**) reads the selected file twice: once as a **blob URL** for the **`<audio>`** element, and once via **`FileReader.readAsArrayBuffer`** for metadata (capped at **40 MiB**).

- **`audiometa.ParseTags`** chooses a parser by signature: **FLAC** (`fLaC` + Vorbis comment block), **MP4** (first top-level **`ftyp`**, then a box walk for **`ilst`**), otherwise **MP3** (**ID3** / **ID3v1**).
- **FLAC**: **`TITLE`**, **`ARTIST`**, **`ALBUM`** from the Vorbis comment block.
- **M4A / MP4** (includes **AAC** and **ALAC** in the same container): iTunes **`ilst`** — **`©nam`**, **`©ART`**, **`©alb`**; **`aART`** fills **artist** only if **`©ART`** is empty.
- **MP3**: **TIT2**, **TPE1**, **TALB** and legacy **ID3v1**.

The UI shows **`Artist — Title · Album`** when any field is present; otherwise only **`:: filename`**. Exotic atoms (**`----`**, odd **data** encodings) may still decode empty.

---

## ZClock and UTC offset

The settings panel stores an integer **UTC offset in hours** added to **UTC** when rendering the clock (see **`apps/zclock.go`**).

- Valid range is roughly **−12 … +14** (what **`+`** / **`−`** adjust).

---

## Guest manifest and runtimes (`entry/guest/<dir>/guest.manifest.json`)

Guestgen emits a sorted list consumed by **`assets/kernel.js`** as **`globalThis.__RIWO_GENERATED_MANIFEST.guests`**. Each row includes **`launchName`** (launcher label), **`runtime`**, and optional **`jsModule`** / **`wasm`**.

| `runtime` | Meaning |
|-----------|---------|
| **`go`** (default when **`main.go`** exists) | Build **`build/<dirname>.wasm`** via **`make guest-wasm`**. **`RunGuestApp`** in the guest uses **`__riwoKernel.consumeGuestBootstrap()`** for pane + **`windowId`**. **`wasm`** in the manifest overrides the relative URL (default **`build/<dir>.wasm`**). |
| **`js-mount`** | No Go **`main.go`**. **`jsModule`** is a relative URL imported with **`import()`**; module must **`export`** **`mount`** (or a default function). **`mount(hostElement, { windowId, launchName, kernel })`** runs in the content pane—**no iframe**. |
| **`wasm-raw`** | No **`main.go`**. **`wasm`** points at a wasm file; kernel uses **`WebAssembly.instantiateStreaming`** with an **empty** import object and calls **`_start`** or **`main`** if present—**best effort** only; WASI-heavy modules need richer imports later. |

If **`guest.manifest.json`** is missing, a directory **with** **`main.go`** is still inferred as **`go`**. **`js-mount`** without explicit **`runtime`** is inferred when **`jsModule`** is set and there is **no** **`main.go`**. Directories that do not qualify are skipped entirely.

Duplicate **`launchName`** values across guests are rejected by guestgen.

---

## Build

```bash
make default       # guestgen stamp → wm.wasm + all Go guests into build/
make wm-wasm       # guestgen + build/wm.wasm only
make guest-wasm    # guestgen + Go guests listed in build/go-guest-dirs.txt
make tiny          # guestgen + TinyGo wm + Go guests (uses tinygo_wasm_exec.js, not go_wasm_exec)
make clean         # removes entire build/
```

There is **no** hand-maintained **`GUESTS`** list: **`build/go-guest-dirs.txt`** is produced alongside **`assets/generated-manifest.js`**. **`$(GUESTGEN_STAMP)`** reruns **`go run ./tools/guestgen`** when **`tools/guestgen/main.go`** or any **`entry/guest/*/main.go`** or **`guest.manifest.json`** changes.

Optional **`wasm-opt`** is applied when present (`-Oz`; WM Go build also uses **`--enable-bulk-memory-opt`**).

---

## What loads in the browser

1. **`index.html`** — **`build/wasm_exec.js`**, then **`assets/generated-manifest.js`** (defines **`__RIWO_GENERATED_MANIFEST`**), then deferred **`assets/kernel.js`**, **`assets/adaptTouch.js`**.
2. **`kernel.js`** defines **`globalThis.__riwoKernel`**, then starts **`build/wm.wasm`** with **`WebAssembly.instantiateStreaming`** + **`go.run`**.
3. **`LaunchDefault(integer windowId)`** completes a **New** rectangle; **`wm.OpenLaunchpad`** fills that window. **`listGuestApps()`** reads the generated manifest.

Paths are resolved with **`new URL(..., document.baseURI)`**. Deployment must keep **`index.html`**, **`build/`**, **`assets/`**, **`apps/`**, **`apps/Monaco/`**, etc. in consistent positions or adjust base URL.

---

## Three tiers

### 1. JS kernel (`assets/kernel.js`)

Responsible for:

- Reading **`__RIWO_GENERATED_MANIFEST`** (with an empty fallback if the script is missing before first build).
- **`spawnGuestApp(name, windowId, contentHostHTMLElement)`**: serializes work on an internal **`spawnChain`** queue (avoids races on the Go bootstrap latch and overlapping instantiations).
- **`consumeGuestBootstrap()`**: synchronous handoff used by **`wm.RunGuestApp`** for **`go`** guests (replaces older single-global bootstrap). Includes a watchdog if Go never consumes the latch.
- **`disposeGuestForWindow(windowId)`**: drops tracking and clears bridged menus for that window.
- **`listGuestApps()`**: sorted **`launchName`** list for the launcher.
- **Guest menu bridge**: **`guestContextMenuAppend`**, **`guestContextMenuTitles`**, **`guestContextMenuInvoke`** (callbacks are typically **`syscall/js.Func`** wrappers; **`invokeGoWasmCallback`** handles plain functions or **`.Invoke`**).
- **`js-mount`** guests receive **`minimalKernelAPI()`** on **`ctx.kernel`** (**`resolveURL`**, **`disposeGuestForWindow`**)—extend carefully if shared concerns grow.

### 2. WM (`build/wm.wasm`)

Same as before: desktop windows (frame + content), launcher calling **`spawnGuestApp`**, **Exit** via **`disposeGuestForWindow`** + **`OpenLaunchpad`**, global context menus, gestures, hide list, etc.

### 3. Guests

- **Go**: **`wm.RunGuestApp`** calls **`consumeGuestBootstrap`**, wraps **`pane`** as **`RiwoWindow.Content`**, runs **`apps`** constructor, then **`select {}`** (main does not exit under wasm).
- **`js-mount`**: no Go guest binary; Monaco is **`assets/js-mounts/monaco-mount.js`** configuring AMD **`vs`** path and attaching the editor to **`host`**.
- **`wasm-raw`**: no Go bridge unless you extend the embedder; pane is cleared and **`_start`/`main`** run if exported.

Guests that need wm-visible menus across instances still use **`wm.RegisterGuestContextMenus`** (**`guest_context_menu.go`**) with the kernel bridge.

---

## Window model (**`RiwoWindow`**)

- **`Frame`**: positioned box on the canvas, **`id`**, move/resize/focus.
- **`Content`**: subtree guests should mutate; **`KernelSpawnGuest`** passes **`window.Content.DOM()`** as the host element.

Deleting a window: **`removeWindow`** removes **`Frame`** and calls **`disposeGuestForWindow`**.

---

## Context menu composition

When the wm menu builds for the focused window (see **`wm/api_ctx_menu.go`**):

1. Fixed rows (Move, New, Resize, Delete, Hide, …).
2. **`RiwoWindow.MenuEntries`** (**Exit**, when spawned from launcher).
3. Kernel-backed titles and **`guestContextMenuInvoke(wid, index)`** for guest-defined actions.

---

## Adding a new guest application

**Go wasm guest**

1. Add **`apps/mydemo.go`** with **`MyDemoConstruct(*wm.RiwoWindow)`**.
2. Add **`entry/guest/mydemo/main.go`** with **`package main`**, **`func main() { wm.RunGuestApp(apps.MyDemoConstruct) }`** (match existing patterns).
3. Optionally add **`entry/guest/mydemo/guest.manifest.json`** with **`"launchName": "Mydemo"`** (otherwise guestgen derives a title from the directory name).
4. Run **`make default`** so guestgen picks up the new dir and **`go-guest-dirs.txt`** includes **`mydemo`**.

**`js-mount` guest (no iframe)**

1. Add **`assets/js-mounts/mywidget-mount.js`** exporting **`export async function mount(host, ctx) { ... }`**.
2. Add **`entry/guest/mywidget/guest.manifest.json`** with **`"runtime": "js-mount"`**, **`"jsModule": "assets/js-mounts/mywidget-mount.js"`**, **`"launchName": "…"`**. **Do not** add **`main.go`** in that directory.
3. **`make default`** — no wasm artifact for **`mywidget`**, launcher still lists it from the manifest.

**`wasm-raw`**

Provide **`wasm`** URL in the manifest only; omit **`main.go`**. Understand import limitations noted above.

---

## Further polish (ideas)

These are optional follow-ups, not a committed roadmap:

- **Guest lifecycle**: Today, exiting a guest clears the DOM but does not tear down the Go wasm instance; long-term isolation might use workers or explicit shutdown if the runtime allows it.
- **`wasm-raw`**: Empty import maps; WASI or custom imports need design if you ship non-Go wasm clients.
- **More audio metadata**: Ogg / Opus **Vorbis comments**, or richer MP4 (**`----`**) if needed.
- **Accessibility**: Focus order and keyboard paths for wm chrome and guests.
- **Tests**: Host-side tests for **`wm/`** logic that can be factored away from **`syscall/js`**, plus headless wasm smoke tests if you add CI.
- **`ctx.kernel` for `js-mount`**: If multiple JS guests share behaviors (menus, file pickers), expose small stable APIs here instead of branching in **`kernel.js`** per app.

---

## Limitations / notes

- **No hard unload** of a Go wasm instance: swapping apps or **Exit** clears the DOM and kernel tables; orphaned guest goroutines may still run until reload—acceptable unless guests move to Workers with an explicit DOM bridge.
- **`wasm-raw`**: empty imports; many real-world modules require WASI — treat as experimental until imports are plumbed if needed.
- **TinyGo**: **`make tiny`** uses **`go/tinygo_wasm_exec.js`** as **`build/wasm_exec.js`**; do not mix TinyGo-produced wasm with the standard Go shim or vice versa.

---

## RTFM guest

The **RTFM** app (`apps/rtfm.go`) is end-user help inside the desktop. **`entry/wm/main.go`** may still log a short cheatsheet when the wm starts.
