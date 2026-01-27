.PHONY: clean default tiny

BUILD_DIR ?= build
OUTPUT    ?= $(BUILD_DIR)/main.wasm
HAS_WASM_OPT := $(shell command -v wasm-opt 2> /dev/null)

clean:
	rm -rf $(BUILD_DIR)

default: clean
	@mkdir -p $(dir $(OUTPUT))
	@cp go/go_wasm_exec.js $(BUILD_DIR)/wasm_exec.js
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(OUTPUT) .
	@touch $(BUILD_DIR)/.$(shell date +%Y-%m-%d)
	@echo "Built default:"
	@ls -lah $(OUTPUT) | awk '{print $$5}'
ifdef HAS_WASM_OPT
	wasm-opt -Oz --enable-bulk-memory-opt -o $(OUTPUT) $(OUTPUT)
	@ls -lah $(OUTPUT) | awk '{print $$5}'
endif

tiny: clean
	@mkdir -p $(dir $(OUTPUT))
	@cp go/tinygo_wasm_exec.js $(BUILD_DIR)/wasm_exec.js
	tinygo build -o $(OUTPUT) -target=wasm -no-debug .
	@touch $(BUILD_DIR)/.$(shell date +%Y-%m-%d)
	@echo "Built tiny:"
	@ls -lah $(OUTPUT) | awk '{print $$5}'
ifdef HAS_WASM_OPT
	wasm-opt -Oz -o $(OUTPUT) $(OUTPUT)
	@ls -lah $(OUTPUT) | awk '{print $$5}'
endif


