.PHONY: all build test wasm clean

all: build wasm

# Build the native CLI binary.
build:
	go build .

# Run all tests.
test:
	go test ./...

# Build the WebAssembly binary and copy the required wasm_exec.js helper.
# The output files are placed in the wasm/ directory:
#   wasm/crslang.wasm  – the compiled WASM module
#   wasm/wasm_exec.js  – the Go-provided JS glue file required to load the WASM
wasm:
	GOOS=js GOARCH=wasm go build -o wasm/crslang.wasm ./wasm/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/

playground:
	mkdir playground
	cp wasm/wasm_exec.js playground/
	cp wasm/crslang.wasm playground/
	cp wasm/index.html playground/

# Remove build artefacts.
clean:
	rm -f crslang wasm/crslang.wasm wasm/wasm_exec.js
