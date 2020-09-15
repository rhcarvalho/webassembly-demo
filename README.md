# webassembly-demo

Experiment with Go, TinyGo and WebAssembly.

https://github.com/golang/go/wiki/WebAssembly
https://tinygo.org/webassembly/webassembly/

## Building

Using the standard Go compiler:

```
GOOS=js GOARCH=wasm go build -o demo.wasm
```

Using TinyGo:

```
tinygo build -o demo-tiny.wasm -target wasm
```

## Running

Random free port:

```
go run .
```

Chosen port:

```
go run . -p 8080
```

Chosen port and WebAssembly file compiled with TinyGo:

```
go run . -p 8080 -wasm=demo-tiny.wasm -tinygo
```
