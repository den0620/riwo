(function (global){
  'use strict';
  global.__RIWO_GENERATED_MANIFEST = {
    version: 2,
    guests: [
	{
		"dir": "dplayer",
		"launchName": "DPlayer",
		"runtime": "go",
		"wasm": "build/dplayer.wasm"
	},
	{
		"dir": "mahjongg",
		"launchName": "Mahjongg",
		"runtime": "go",
		"wasm": "build/mahjongg.wasm"
	},
	{
		"dir": "monaco",
		"launchName": "Monaco",
		"runtime": "js-mount",
		"jsModule": "assets/js-mounts/monaco-mount.js"
	},
	{
		"dir": "rtfm",
		"launchName": "RTFM",
		"runtime": "go",
		"wasm": "build/rtfm.wasm"
	},
	{
		"dir": "zclock",
		"launchName": "ZClock",
		"runtime": "go",
		"wasm": "build/zclock.wasm"
	}
]
  };
})(typeof globalThis !== 'undefined' ? globalThis : this);
