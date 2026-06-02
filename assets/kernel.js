'use strict';

(function () {
	/** Fallback if generated manifest not loaded yet (run `make` / guestgen first). */
	const FALLBACK_GUESTS = [];

	const guestsTracking = new Map();

	/** @typedef {{ name: string, fn: any }} GuestMenuRow */
	/** @type {Map<number, GuestMenuRow[]>} */
	const guestContextMenusByWindow = new Map();

	function guestContextMenuClear(windowId) {
		guestContextMenusByWindow.delete(Number(windowId));
	}

	function guestContextMenuAppend(windowId, name, goCb) {
		const wid = Number(windowId);
		if (!guestContextMenusByWindow.has(wid)) {
			guestContextMenusByWindow.set(wid, []);
		}
		guestContextMenusByWindow.get(wid).push({ name: String(name), fn: goCb });
	}

	function guestContextMenuTitles(windowId) {
		const rows = guestContextMenusByWindow.get(Number(windowId)) || [];
		return rows.map((r) => r.name);
	}

	function invokeGoWasmCallback(fn) {
		if (typeof fn === 'function') {
			return fn();
		}
		if (fn != null && typeof fn.Invoke === 'function') {
			return fn.Invoke();
		}
		return undefined;
	}

	function guestContextMenuInvoke(windowId, index) {
		const rows = guestContextMenusByWindow.get(Number(windowId));
		if (!rows) {
			return;
		}
		const i = Number(index);
		if (i < 0 || i >= rows.length) {
			return;
		}
		invokeGoWasmCallback(rows[i].fn);
	}

	function disposeGuestForWindow(windowId) {
		guestsTracking.delete(Number(windowId));
		guestContextMenuClear(windowId);
	}

	function resolveURL(path) {
		return new URL(path, document.baseURI).href;
	}

	function readGuestManifestRecords() {
		const m = globalThis.__RIWO_GENERATED_MANIFEST;
		const list = Array.isArray(m?.guests) ? m.guests : FALLBACK_GUESTS;
		return list;
	}

	function manifestByLaunchName() {
		/** @type {Map<string, any>} */
		const map = new Map();
		for (const row of readGuestManifestRecords()) {
			if (!row.launchName || !row.runtime) {
				continue;
			}
			map.set(String(row.launchName), row);
		}
		return map;
	}

	function listGuestApps() {
		return readGuestManifestRecords()
			.map((r) => r.launchName)
			.filter(Boolean)
			.sort((a, b) => String(a).localeCompare(String(b), undefined, { sensitivity: 'base' }));
	}

	/** One guest spawn at a time; avoids races on latch + Go startup. */
	let spawnChain = Promise.resolve();

	function enqueueSpawn(fn) {
		spawnChain = spawnChain.then(fn).catch((e) => {
			console.error('[riwo kernel] spawn error:', e);
		});
		return spawnChain;
	}

	/** Pending Go guest reads bootstrap via consumeGuestBootstrap (sync Go main startup). */
	let guestBootstrapThunk = null;
	let guestBootstrapResolve = null;
	let goBootstrapWatchdog = null;

	function consumeGuestBootstrap() {
		if (typeof guestBootstrapThunk !== 'function') {
			return null;
		}
		const out = guestBootstrapThunk();
		guestBootstrapThunk = null;
		if (goBootstrapWatchdog !== null) {
			clearTimeout(goBootstrapWatchdog);
			goBootstrapWatchdog = null;
		}
		if (typeof guestBootstrapResolve === 'function') {
			const r = guestBootstrapResolve;
			guestBootstrapResolve = null;
			r();
		}
		return out;
	}

	async function runGoGuest(manifestRow, wid, pane) {
		disposeGuestForWindow(wid);
		pane.replaceChildren();

		await new Promise((resolveBootstrap, rejectBootstrap) => {
			guestBootstrapThunk = () => ({
				windowId: Number(wid),
				pane
			});
			guestBootstrapResolve = resolveBootstrap;

			goBootstrapWatchdog = setTimeout(() => {
				if (typeof guestBootstrapThunk === 'function') {
					console.error(
						'[riwo kernel] Go guest did not call consumeGuestBootstrap (RunGuestApp must run first); unblocking spawn queue'
					);
					guestBootstrapThunk = null;
					guestBootstrapResolve = null;
					goBootstrapWatchdog = null;
					rejectBootstrap(new Error('consumeGuestBootstrap watchdog'));
				}
			}, 15000);

			void (async () => {
				try {
					const go = new Go();
					const res = await WebAssembly.instantiateStreaming(
						fetch(resolveURL(manifestRow.wasm)),
						go.importObject
					);
					guestsTracking.set(Number(wid), {
						wasmPath: manifestRow.wasm,
						runtime: 'go',
						windowId: wid
					});
					void go.run(res.instance);
				} catch (e) {
					if (goBootstrapWatchdog !== null) {
						clearTimeout(goBootstrapWatchdog);
						goBootstrapWatchdog = null;
					}
					guestBootstrapThunk = null;
					guestBootstrapResolve = null;
					rejectBootstrap(e);
				}
			})();
		});
	}

	async function runJsMountGuest(jsModuleRel, manifestRow, wid, pane) {
		disposeGuestForWindow(wid);
		pane.replaceChildren();

		const mod = await import(resolveURL(jsModuleRel));
		const fn = typeof mod.mount === 'function' ? mod.mount : typeof mod.default === 'function' ? mod.default : null;

		if (typeof fn !== 'function') {
			throw new Error('JS mount module exports no mount() — require export async function mount(host, ctx)');
		}

		await fn(pane, { windowId: Number(wid), launchName: manifestRow.launchName, kernel: minimalKernelAPI() });
		guestsTracking.set(Number(wid), { runtime: 'js-mount', windowId: wid, module: jsModuleRel });
	}

	function minimalKernelAPI() {
		return {
			resolveURL,
			disposeGuestForWindow
		};
	}

	async function runRawWasmGuest(wasmPath, wid, pane) {
		disposeGuestForWindow(wid);
		pane.replaceChildren();

		const res = await WebAssembly.instantiateStreaming(fetch(resolveURL(wasmPath)), {});
		const ex = res.instance.exports;
		if (typeof ex._start === 'function') {
			ex._start();
		} else if (typeof ex.main === 'function') {
			ex.main();
		}
		guestsTracking.set(Number(wid), { runtime: 'wasm-raw', windowId: wid, wasmPath });
	}

	function spawnGuestApp(appName, windowId, contentHostJsValue) {
		if (!contentHostJsValue || typeof contentHostJsValue.replaceChildren !== 'function') {
			console.warn('[riwo kernel] Invalid content host element');
			return Promise.resolve();
		}
		const byName = manifestByLaunchName();
		const row = byName.get(String(appName));
		if (!row) {
			console.warn('[riwo kernel] Unknown guest:', appName);
			return Promise.resolve();
		}

		return enqueueSpawn(async () => {
			const wid = Number(windowId);
			const pane = contentHostJsValue;

			switch (row.runtime) {
				case 'go':
					await runGoGuest(row, wid, pane);
					return;
				case 'js-mount':
					await runJsMountGuest(row.jsModule, row, wid, pane);
					return;
				case 'wasm-raw':
					await runRawWasmGuest(row.wasm, wid, pane);
					return;
				default:
					console.warn('[riwo kernel] Unsupported runtime', row.runtime);
			}
		});
	}

	async function startWM(wasmRel) {
		const go = new Go();
		try {
			const res = await WebAssembly.instantiateStreaming(fetch(resolveURL(wasmRel)), go.importObject);
			go.run(res.instance);
		} catch (e) {
			console.error('[riwo kernel] wm failed to start:', e);
		}
	}

	globalThis.__riwoKernel = {
		spawnGuestApp,
		disposeGuestForWindow,
		listGuestApps,
		guestContextMenuAppend,
		guestContextMenuTitles,
		guestContextMenuInvoke,
		consumeGuestBootstrap,
		startWM
	};

	function kickStartWM() {
		startWM('build/wm.wasm');
	}
	if (document.readyState === 'loading') {
		window.addEventListener('DOMContentLoaded', kickStartWM);
	} else {
		queueMicrotask(kickStartWM);
	}
})();
