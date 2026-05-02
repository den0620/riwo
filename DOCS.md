# Riwo

Riwo is a Rio-inspired window shell in the browser: the **window manager** and each **application** compile to separate **WebAssembly** modules. A small **JavaScript kernel** owns page load paths, chooses which wasm binary to instantiate, holds cross-instance bookkeeping (guest context menus), and hosts **Monaco** without an iframe.

This document reflects the architecture after the split from the earlier single `main.wasm` tree.

---

## Repository layout

| Path | Role |
|------|------|
| **`entry/wm`** | WM `main`: registers globals (`Logging`, `LaunchDefault`), starts `wm.InitializeContextMenu`, `InitializeGlobalMouseEvents`. |
| **`entry/guest/<name>`** | One `main` per launcher app (ZClock, DPlayer, Mahjongg, RTFM, Monaco). |
| **`wm/`** | Window manager DOM helpers, gestures, launcher, syscall bridges to **`__riwoKernel`**. |
| **`apps/`** | Guest UI constructors (`*_Construct`), imported only by **`entry/guest/*`**. |
| **`assets/kernel.js`** | Bootstrap: **`__riwoKernel`**, guest spawn, Monaco mount, bridged menus. |
| **`build/`** | Generated: **`wasm_exec.js`**, **`wm.wasm`**, **`<guest>.wasm`**. |
| **`go/go_wasm_exec.js`** | Copied into **`build/wasm_exec.js`** for canonical Go wasm. |
| **`apps/Monaco/`** | Static Monaco assets (AMD loader + `vs/*`), loaded by the kernel. |

Module path: **`module riwo`** (`go.mod`).

---

## Build

```bash
make default       # wm + all guests = build/
make wm-wasm       # only build/wm.wasm
make guest-wasm    # only guests (see GUESTS)
make tiny          # TinyGo wm + guests; replaces wasm_exec with tinygo variant; rm -rf build first
make clean         # removes build/
```

Configure guest list via **`Makefile` `GUESTS`**. Optional **`wasm-opt`** is applied when present (`-Oz`; WM build also uses **`--enable-bulk-memory-opt`**).

---

## What loads in the browser

1. **`index.html`** - **`build/wasm_exec.js`** (Go runtime shim), **`defer`** **`assets/kernel.js`**, **`assets/adaptTouch.js`**.
2. **`kernel.js`** defines **`globalThis.__riwoKernel`**, then starts **`build/wm.wasm`** with **`WebAssembly.instantiateStreaming`** + **`go.run`**.
3. **`LaunchDefault(integer windowId)`** is invoked when a **New** rectangle completes; **`wm.OpenLaunchpad`** fills that window.

Paths under **`kernel.js`** (**`build/wm.wasm`**, **`build/<guest>.wasm`**) are resolved with **`new URL(..., document.baseURI)`**, so deployment must keep **`index.html`**, **`build/`**, **`assets/`**, and **`apps/Monaco/`** in consistent relative positions (or serve from the repo root).

---

## Three tiers

### 1. JS kernel (`assets/kernel.js`)

Responsible for:

- **`GUEST_WASM`**: launcher name -> relative wasm URL (must match **`Makefile`** `GUESTS` and tile labels from **`listGuestApps()`**).
- **`spawnGuestApp(name, windowId, contentHostHTMLElement)`**: `disposeGuestForWindow`, clears the pane, sets **`globalThis.__riwoGuestBootstrap`**, **`new Go()`**, instantiate guest wasm, **`go.run`**.
- **`disposeGuestForWindow(windowId)`**: drops kernel tracking and **clears bridged guest context menus** for that window.
- **`listGuestApps()`**: exposes sorted keys of **`GUEST_WASM`** to the WM via syscall.
- **Guest menu bridge**: **`guestContextMenuAppend`**, **`guestContextMenuTitles`**, **`guestContextMenuInvoke`**. Callbacks registered from guests are **`syscall/js.Func`** values exposed to JS as **plain functions**; **`invokeGoWasmCallback`** calls **`fn()`** (not **`fn.Invoke`**, unless a future shim provides it).
- **Monaco**: **`monacoMount`**, **`monacoMountDone`** (loads **`apps/Monaco/min/vs/loader.js`**, configures `vs` path, creates editor).

### 2. WM (`build/wm.wasm`)

- Builds **desktop windows**: outer **`Frame`** (drag, stack, resize, delete, hide) and inner **`Content`** (launcher tiles or guest root element passed to kernel).
- **Launcher**: **`wm.OpenLaunchpad`** — reads app names via **`__riwoKernel.listGuestApps`**; **`wm.KernelSpawnGuest`** calls **`spawnGuestApp`**. **`MenuEntries`** cleared on launcher.
- **Exit**: **`attachExitMenu`** sets **`RiwoWindow.MenuEntries`** with **「Exit」** calling **`ReturnToLaunchpad`** (kernel **`disposeGuestForWindow`**, then **`OpenLaunchpad`**). Same WM memory as **`api_ctx_menu`**, so no JS menu bridge needed for Exit.
- Exposes **`Logging()`** and **`LaunchDefault(...)`** on **`js.Global`**.
- Implements global context menu, move / resize ghosts, hide list, etc. (see **`wm/`** sources).

### 3. Guests (`build/<name>.wasm`)

- **`wm.RunGuestApp`** reads **`__riwoGuestBootstrap.pane`** and **`windowId`**, wraps the pane as **`RiwoWindow.Content`,** and runs the app constructor from **`riwo/apps`**.
- **`wm.RegisterGuestContextMenus`** publishes per-window rows through the kernel (**`guestContextMenuAppend`**) because guest and WM live in **different wasm instances**; **`MenuEntries` on structs inside guests do not reach the WM.**

Monaco **`entry/guest/monaco`** is minimal: lays out **`Content`** and calls **`monacoMountDone`** with a **`js.Func`** so initialization can finish asynchronously.

---

## Window model (**`RiwoWindow`**)

- **`Frame`**: positioned box on the canvas, carries window **`id`** in the DOM sense, owns move/resize/focus handlers.
- **`Content`**: the only subtree guests and launcher should mutate; **`KernelSpawnGuest`** receives **`window.Content.DOM()`**.

Deleting a window: **`removeWindow`** removes **`Frame`** and calls **`disposeGuestForWindow`**.

---

## Context menu composition

When the wm menu builds for the focused window (see **`wm/api_ctx_menu.go`**):

1. Fixed rows (Move, New, Resize, Delete, Hide, ...).
2. **`RiwoWindow.MenuEntries`** (**Exit**, when a guest was started from the launcher).
3. Kernel-backed titles followed by **`guestContextMenuInvoke(wid, index)`** for guest-defined actions.

Guests use **`wm.RegisterGuestContextMenus`** (**`wm/guest_context_menu.go`**).

---

## Adding a new guest application

1. Add **`apps/mydemo.go`** (or similar) with an exported **`MyDemoConstruct(*wm.RiwoWindow)`** (or **`RunGuestApp`** body).
2. Add **`entry/guest/mydemo/main.go`** with **`package main`**, **`func main() { wm.RunGuestApp(apps.MyDemoConstruct) }`** (pattern match existing guests).
3. Append **`mydemo`** to **`GUESTS`** in the **`Makefile`**.
4. Add **`Mydemo: 'build/mydemo.wasm'`** (matching the **`GUESTS`** make variable and tile label) to **`GUEST_WASM`** in **`assets/kernel.js`**.

Rebuild **`make default`**.

---

## Limitations / notes

- **No hard unload** of a running Go wasm instance: swapping apps or **Exit** clears the DOM and kernel tables; orphaned guest goroutines may still run until tab reload - acceptable tradeoff unless you isolate guests in Workers with an explicit IPC DOM bridge.
- **Monaco**: runs as JS/editor inside the wm content tree; Monaco assets must stay reachable at **`apps/Monaco/`** paths used in **`kernel.js`**.
- **TinyGo**: **`make tiny`** uses **`go/tinygo_wasm_exec.js`** as **`build/wasm_exec.js`**; do not mix TinyGo-produced wasm with the standard Go shim or vice versa.

---

## RTFM guest

The **RTFM** app (`apps/rtfm.go`) is intended as end-user-readable help inside the desktop. For scripting tips, **`entry/wm/main.go`** still logs a short cheatsheet when the wm starts.
