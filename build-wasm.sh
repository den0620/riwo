#!/bin/env bash
GOOS=js GOARCH=wasm go build -o ./build/main.wasm .
touch ./build/.$(printf "%(%Y-%m-%d)T\n" -1)

