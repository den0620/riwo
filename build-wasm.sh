#!/bin/env bash
cd src
GOOS=js GOARCH=wasm go build -o ../build/main.wasm .
cd ..
touch ./build/.$(printf "%(%Y-%m-%d)T\n" -1)

