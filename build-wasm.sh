#!/bin/env bash
GOOS=js GOARCH=wasm go build -o ./build/main.wasm ./src/main.go
touch ./build/.$(printf "%(%Y-%m-%d)T\n" -1)

