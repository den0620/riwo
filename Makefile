# Standard Go wasm: builds build/wm.wasm plus build/<guest>.wasm (see GUESTS).
# TinyGo wasm: same outputs with `make tiny` (needs tinygo installed; wipes build/ first).

.PHONY: clean guest-wasm wm-wasm all default tiny

BUILD_DIR ?= build
OUTPUT_WM  ?= $(BUILD_DIR)/wm.wasm
LDFLAGS    ?= -s -w
GUESTS     ?= zclock dplayer mahjongg rtfm monaco

HAS_WASM_OPT := $(shell command -v wasm-opt 2> /dev/null)

clean:
	rm -rf $(BUILD_DIR)

wm-wasm: $(BUILD_DIR)/wasm_exec.js
	@mkdir -p $(dir $(OUTPUT_WM))
	GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_WM) ./entry/wm
ifdef HAS_WASM_OPT
	wasm-opt -Oz --enable-bulk-memory-opt -o $(OUTPUT_WM) $(OUTPUT_WM)
endif
	@echo "$(OUTPUT_WM):" && ls -lah $(OUTPUT_WM) | awk '{print $$5}'

guest-wasm: $(BUILD_DIR)/wasm_exec.js
	@mkdir -p $(BUILD_DIR)
	@for g in $(GUESTS); do \
		out="$(BUILD_DIR)/$$g.wasm"; \
		echo "guest -> $$out"; \
		GOOS=js GOARCH=wasm go build -ldflags="$(LDFLAGS)" -o "$$out" "./entry/guest/$$g" || exit 1; \
		ls -lah "$$out" | awk '{print $$5}'; \
		if command -v wasm-opt >/dev/null 2>&1; then wasm-opt -Oz --enable-bulk-memory-opt -o "$$out" "$$out"; ls -lah "$$out" | awk '{print $$5}'; fi; \
	done

$(BUILD_DIR)/wasm_exec.js: go/go_wasm_exec.js
	@mkdir -p $(BUILD_DIR)
	@cp go/go_wasm_exec.js $(BUILD_DIR)/wasm_exec.js

default: wm-wasm guest-wasm
	@echo "Built WM + guests into $(BUILD_DIR)/"

tiny: clean
	@mkdir -p $(dir $(OUTPUT_WM))
	@cp go/tinygo_wasm_exec.js $(BUILD_DIR)/wasm_exec.js
	tinygo build -o $(OUTPUT_WM) -target=wasm -no-debug ./entry/wm
	@for g in $(GUESTS); do \
		out="$(BUILD_DIR)/$$g.wasm"; \
		tinygo build -o "$$out" -target=wasm -no-debug "./entry/guest/$$g" || exit 1; \
	done
ifdef HAS_WASM_OPT
	wasm-opt -Oz -o $(OUTPUT_WM) $(OUTPUT_WM)
	@for g in $(GUESTS); do wasm-opt -Oz -o "$(BUILD_DIR)/$$g.wasm" "$(BUILD_DIR)/$$g.wasm"; done
endif

all: default
