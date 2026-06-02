# Standard Go wasm: builds build/wm.wasm plus Go guests listed in build/go-guest-dirs.txt (from tools/guestgen).
# Launcher manifest: assets/generated-manifest.js
# TinyGo: `make tiny` (runs guestgen fresh then tinygo loops over go-guest-dirs.txt).

.PHONY: clean guest-wasm wm-wasm all default tiny

BUILD_DIR ?= build
OUTPUT_WM  ?= $(BUILD_DIR)/wm.wasm
LDFLAGS    ?= -s -w

GUESTGEN_STAMP := $(BUILD_DIR)/.guestgen_stamp

HAS_WASM_OPT := $(shell command -v wasm-opt 2> /dev/null)

$(GUESTGEN_STAMP): tools/guestgen/main.go $(wildcard entry/guest/*/main.go entry/guest/*/guest.manifest.json)
	@mkdir -p $(BUILD_DIR) assets
	go run ./tools/guestgen -entry entry/guest -outjs assets/generated-manifest.js
	@touch $@

clean:
	rm -rf $(BUILD_DIR)

wm-wasm: $(GUESTGEN_STAMP) $(BUILD_DIR)/wasm_exec.js
	@mkdir -p $(dir $(OUTPUT_WM))
	GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_WM) ./entry/wm
ifdef HAS_WASM_OPT
	wasm-opt -Oz --enable-bulk-memory-opt -o $(OUTPUT_WM) $(OUTPUT_WM)
endif
	@echo "$(OUTPUT_WM):" && ls -lah $(OUTPUT_WM) | awk '{print $$5}'

guest-wasm: $(GUESTGEN_STAMP) $(BUILD_DIR)/wasm_exec.js
	@mkdir -p $(BUILD_DIR)
	@if [ ! -s $(BUILD_DIR)/go-guest-dirs.txt ]; then echo "guest-wasm: no Go wasm guests (js-mount only?)."; fi
	@while IFS= read -r g; do \
		test -n "$$g" || continue; \
		out="$(BUILD_DIR)/$$g.wasm"; \
		echo "guest -> $$out"; \
		GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -o "$$out" "./entry/guest/$$g" || exit 1; \
		ls -lah "$$out" | awk '{print $$5}'; \
		if command -v wasm-opt >/dev/null 2>&1; then wasm-opt -Oz --enable-bulk-memory-opt -o "$$out" "$$out"; ls -lah "$$out" | awk '{print $$5}'; fi; \
	done < $(BUILD_DIR)/go-guest-dirs.txt

$(BUILD_DIR)/wasm_exec.js: go/go_wasm_exec.js
	@mkdir -p $(BUILD_DIR)
	@cp go/go_wasm_exec.js $(BUILD_DIR)/wasm_exec.js

default: wm-wasm guest-wasm
	@echo "Built WM + guests into $(BUILD_DIR)/"

tiny:
	@mkdir -p $(BUILD_DIR) assets
	go run ./tools/guestgen -entry entry/guest -outjs assets/generated-manifest.js
	@mkdir -p $(dir $(OUTPUT_WM))
	@cp go/tinygo_wasm_exec.js $(BUILD_DIR)/wasm_exec.js
	tinygo build -o $(OUTPUT_WM) -target=wasm -no-debug ./entry/wm
	@while IFS= read -r g; do \
		test -n "$$g" || continue; \
		out="$(BUILD_DIR)/$$g.wasm"; \
		echo "tiny guest -> $$out"; \
		tinygo build -o "$$out" -target=wasm -no-debug "./entry/guest/$$g" || exit 1; \
	done < $(BUILD_DIR)/go-guest-dirs.txt
ifdef HAS_WASM_OPT
	wasm-opt -Oz -o $(OUTPUT_WM) $(OUTPUT_WM)
	@while IFS= read -r g; do \
		test -n "$$g" || continue; \
		wasm-opt -Oz -o "$(BUILD_DIR)/$$g.wasm" "$(BUILD_DIR)/$$g.wasm"; \
	done < $(BUILD_DIR)/go-guest-dirs.txt
endif

all: default
