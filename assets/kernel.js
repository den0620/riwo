'use strict';

(function () {
	/** Wasm paths resolved relative to the page (same folder as index.html by default). */
	const GUEST_WASM = Object.freeze({
		ZClock: 'build/zclock.wasm',
		DPlayer: 'build/dplayer.wasm',
		Mahjongg: 'build/mahjongg.wasm',
		RTFM: 'build/rtfm.wasm',
		Monaco: 'build/monaco.wasm'
	});

	/** Best-effort tracking of guest goroutines keyed by wm window id. */
	const guests = new Map();

	/** @typedef {{ name: string, fn: any }} GuestMenuRow */
	/** Guest context menu callbacks live in JS so WM wasm can trigger them cross-instance. */
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

	/** Go syscall/js.Func appears in JS as a plain function (wasm_exec `_makeFuncWrapper`), not `{ Invoke }`. */
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
		guests.delete(Number(windowId));
		guestContextMenuClear(windowId);
	}

	function resolveURL(path) {
		return new URL(path, document.baseURI).href;
	}

	function listGuestApps() {
		return Object.keys(GUEST_WASM).slice().sort();
	}

	function loaderScriptLoaded(src) {
		return Array.from(document.scripts).some((s) => s.src === src);
	}

	function loadScriptOnce(src) {
		const abs = resolveURL(src);
		if (loaderScriptLoaded(abs)) {
			return Promise.resolve();
		}
		return new Promise((resolve, reject) => {
			const script = document.createElement('script');
			script.src = abs;
			script.onload = () => resolve();
			script.onerror = () => reject(new Error('Failed to load ' + abs));
			document.head.appendChild(script);
		});
	}

	/**
	 * Mount Monaco Editor into hostEl (fills host). Loads AMD loader from apps/Monaco once.
	 */
	async function monacoMount(hostEl) {
		hostEl.replaceChildren();

		const wrap = document.createElement('div');
		wrap.style.width = '100%';
		wrap.style.height = '100%';
		wrap.style.margin = '0';
		wrap.style.padding = '0';
		wrap.style.overflow = 'hidden';

		const container = document.createElement('div');
		container.style.width = '100%';
		container.style.height = '100%';

		const editor = document.createElement('div');
		editor.id = 'riwo-monaco-' + Date.now().toString(36);
		editor.style.width = '100%';
		editor.style.height = '100%';

		container.appendChild(editor);
		wrap.appendChild(container);
		hostEl.appendChild(wrap);

		await loadScriptOnce('apps/Monaco/min/vs/loader.js');

		await new Promise((resolve, reject) => {
			if (typeof window.require !== 'function' || typeof window.require.config !== 'function') {
				reject(new Error('Monaco loader (require) is not available'));
				return;
			}
			const vsHref = resolveURL('apps/Monaco/min/vs');

			window.require.config({
				paths: { vs: vsHref },
				catchError: true
			});

			window.require(['vs/editor/editor.main'], function () {
				try {
					const dark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
					window.monaco.editor.create(editor, {
						value: '// Ready\n',
						language: 'go',
						theme: dark ? 'vs-dark' : 'vs-light',
						automaticLayout: true,
						minimap: { enabled: false },
						scrollBeyondLastLine: false
					});
				} catch (e) {
					reject(e);
					return;
				}
				resolve();
			});
		});
	}

	/**
	 * Same as monacoMount but invokes onDone JS function once editor is ready (or on error — still invokes).
	 */
	function monacoMountDone(hostEl, onDoneJs) {
		monacoMount(hostEl).then(
			() => {
				invokeGoWasmCallback(onDoneJs);
			},
			(err) => {
				console.error(err);
				invokeGoWasmCallback(onDoneJs);
			}
		);
	}

	async function spawnGuestApp(appName, windowId, contentHostJsValue) {
		const wasmRel = GUEST_WASM[appName];
		if (!wasmRel) {
			console.warn('[riwo kernel] Unknown guest:', appName);
			return;
		}
		if (!contentHostJsValue || typeof contentHostJsValue.replaceChildren !== 'function') {
			console.warn('[riwo kernel] Invalid content host element');
			return;
		}
		disposeGuestForWindow(windowId);
		contentHostJsValue.replaceChildren();

		globalThis.__riwoGuestBootstrap = {
			windowId: Number(windowId),
			pane: contentHostJsValue
		};

		try {
			const go = new Go();
			const res = await WebAssembly.instantiateStreaming(
				fetch(resolveURL(wasmRel)),
				go.importObject
			);
			guests.set(Number(windowId), { wasmPath: wasmRel, windowId });
			go.run(res.instance);
		} catch (e) {
			console.error('[riwo kernel] spawnGuestApp:', appName, e);
			guests.delete(Number(windowId));
		}
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
		monacoMount,
		monacoMountDone,
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
