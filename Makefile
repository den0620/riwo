.PHONY: all clean

BUILD_DIR ?= build
OUTPUT    ?= $(BUILD_DIR)/main.wasm

all:
	@mkdir -p $(dir $(OUTPUT))
	GOOS=js GOARCH=wasm go build -o $(OUTPUT) .
	@touch $(BUILD_DIR)/.$(shell date +%Y-%m-%d)

clean:
	rm -r $(BUILD_DIR)

