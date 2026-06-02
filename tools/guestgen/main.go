/*
Guestgen discovers entry/guest/<dir>/ and emits:
  - build/go-guest-dirs.txt   directory names with Go main packages (one per line, for make loops)
  - assets/generated-manifest.js  global __RIWO_GENERATED_MANIFEST for kernel.js

Per-directory rules:
  - guest.manifest.json (optional overrides)
  - else if main.go exists -> runtime go, wasm build/<dirname>.wasm
  - runtime js-mount requires guest.manifest.json with jsModule; no wasm build line
  - runtime wasm-raw requires wasm path and no Go main.go
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type guestManifestFile struct {
	LaunchName string `json:"launchName"`
	Runtime    string `json:"runtime"`
	JSModule   string `json:"jsModule"`
	WASM       string `json:"wasm"`
}

type guestRecord struct {
	Dir        string `json:"dir"`
	LaunchName string `json:"launchName"`
	Runtime    string `json:"runtime"` // go | js-mount | wasm-raw
	JSModule   string `json:"jsModule,omitempty"`
	WASM       string `json:"wasm,omitempty"`
}

func main() {
	entryRoot := flag.String("entry", "entry/guest", "Directory containing guest subfolders")
	outJS := flag.String("outjs", "assets/generated-manifest.js", "JS manifest consumed by kernel.js")
	flag.Parse()

	if err := run(*entryRoot, *outJS); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(entryRoot, outJS string) error {
	ents, err := os.ReadDir(entryRoot)
	if err != nil {
		return fmt.Errorf("read %s: %w", entryRoot, err)
	}

	var guests []guestRecord
	for _, ent := range ents {
		if !ent.IsDir() || strings.HasPrefix(ent.Name(), ".") {
			continue
		}
		dirName := ent.Name()
		base := filepath.Join(entryRoot, dirName)
		manPath := filepath.Join(base, "guest.manifest.json")
		mainPath := filepath.Join(base, "main.go")

		hasMain := false
		if _, stat := os.Stat(mainPath); stat == nil {
			hasMain = true
		}

		var gf guestManifestFile
		if dat, readErr := os.ReadFile(manPath); readErr == nil {
			if err := json.Unmarshal(dat, &gf); err != nil {
				return fmt.Errorf("%s: invalid JSON: %w", manPath, err)
			}
		}

		rt := strings.TrimSpace(strings.ToLower(gf.Runtime))
		if rt == "" && gf.JSModule != "" && !hasMain {
			rt = "js-mount"
		}
		if rt == "" && hasMain {
			rt = "go"
		}
		if rt == "" {
			continue
		}

		if rt == "js-mount" && hasMain {
			return fmt.Errorf("%s: js-mount cannot include main.go; use wasm guest or plain JS mount module", dirName)
		}

		display := strings.TrimSpace(gf.LaunchName)
		if display == "" {
			display = asciiTitle(dirName)
		}

		var wasmRel string
		switch rt {
		case "go":
			if !hasMain {
				return fmt.Errorf("%s: runtime go requires main.go", dirName)
			}
			wasmRel = strings.TrimSpace(gf.WASM)
			if wasmRel == "" {
				wasmRel = filepath.ToSlash(filepath.Join("build", dirName+".wasm"))
			}
		case "js-mount":
			if gf.JSModule == "" {
				return fmt.Errorf("%s: js-mount requires jsModule in guest.manifest.json", dirName)
			}
		case "wasm-raw":
			wasmRel = strings.TrimSpace(gf.WASM)
			if wasmRel == "" {
				return fmt.Errorf("%s: wasm-raw requires \"wasm\" path in manifest", dirName)
			}
			if hasMain {
				return fmt.Errorf("%s: wasm-raw should not bundle main.go (use wasm path only)", dirName)
			}
		default:
			return fmt.Errorf("%s: unsupported runtime %q", dirName, rt)
		}

		guests = append(guests, guestRecord{
			Dir:        dirName,
			LaunchName: display,
			Runtime:    rt,
			JSModule:   strings.TrimSpace(gf.JSModule),
			WASM:       wasmRel,
		})
	}

	if len(guests) == 0 {
		return fmt.Errorf("no guests under %s", entryRoot)
	}

	sort.Slice(guests, func(i, j int) bool {
		return strings.ToLower(guests[i].LaunchName) < strings.ToLower(guests[j].LaunchName)
	})

	launchSeen := map[string]string{}
	var goGuests []string
	for _, g := range guests {
		if prev, dup := launchSeen[g.LaunchName]; dup && prev != g.Dir {
			return fmt.Errorf("duplicate launch name %q in %s and %s", g.LaunchName, prev, g.Dir)
		}
		launchSeen[g.LaunchName] = g.Dir
		if g.Runtime == "go" {
			goGuests = append(goGuests, g.Dir)
		}
	}

	sort.Strings(goGuests)

	if err := writeGoGuestDirs(filepath.Join("build", "go-guest-dirs.txt"), goGuests); err != nil {
		return err
	}
	if err := writeJSManifest(outJS, guests); err != nil {
		return err
	}
	return nil
}

func writeGoGuestDirs(out string, names []string) error {
	dir := filepath.Dir(out)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(n)
		sb.WriteByte('\n')
	}
	return os.WriteFile(out, []byte(sb.String()), 0o644)
}

func writeJSManifest(path string, guests []guestRecord) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(guests, "", "\t")
	if err != nil {
		return err
	}
	payload := `(function (global){
  'use strict';
  global.__RIWO_GENERATED_MANIFEST = {
    version: 2,
    guests: ` + string(data) + `
  };
})(typeof globalThis !== 'undefined' ? globalThis : this);
`
	return os.WriteFile(path, []byte(payload), 0o644)
}

func asciiTitle(s string) string {
	if s == "" {
		return s
	}
	l := strings.ToLower(s)
	return strings.ToUpper(l[:1]) + l[1:]
}
