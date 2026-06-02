/**
 * Monaco editor as a plain JS-mount guest — no wasm, no iframe.
 * @param {HTMLElement} hostEl wm Content pane — filled by this mount
 * @param {{ windowId: number }} _ctx Reserved for telemetry / teardown hooks
 */
export async function mount(hostEl, _ctx) {
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

	const editorRoot = document.createElement('div');
	editorRoot.style.width = '100%';
	editorRoot.style.height = '100%';
	editorRoot.id = 'riwo-monaco-' + Date.now().toString(36);

	container.appendChild(editorRoot);
	wrap.appendChild(container);
	hostEl.appendChild(wrap);

	const base = document.baseURI;
	const vsLoader = new URL('apps/Monaco/min/vs/loader.js', base).href;
	const vsHref = new URL('apps/Monaco/min/vs', base).href;

	await loadScriptOnce(vsLoader);

	await new Promise((resolve, reject) => {
		if (typeof window.require !== 'function' || typeof window.require.config !== 'function') {
			reject(new Error('Monaco AMD loader unavailable'));
			return;
		}
		window.require.config({
			paths: { vs: vsHref },
			catchError: true
		});
		window.require(['vs/editor/editor.main'], function () {
			try {
				const dark =
					window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
				window.monaco.editor.create(editorRoot, {
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
			resolve(undefined);
		});
	});
}

function loadScriptOnce(src) {
	const abs = src;
	const exists = Array.from(document.scripts).some((s) => s.src === abs);
	if (exists) return Promise.resolve();
	return new Promise((resolve, reject) => {
		const script = document.createElement('script');
		script.src = abs;
		script.onload = () => resolve();
		script.onerror = () => reject(new Error('Failed to load ' + abs));
		document.head.appendChild(script);
	});
}
