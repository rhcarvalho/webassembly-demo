// +build !wasm

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	port     = flag.Int("p", 0, "port")
	wasmFile = flag.String("wasm", "demo.wasm", "path to compiled WebAssembly file")
	tinygo   = flag.Bool("tinygo", false, "wasm built with tinygo")
)

var wasmExecFile string

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-store") // disable caching
	switch r.URL.Path {
	case "/":
		http.ServeContent(w, r, "index.html", time.Now(), strings.NewReader(`<html>
	<head>
		<meta charset="utf-8"/>
		<script src="wasm_exec.js"></script>
		<script>
			const go = new Go();
			WebAssembly.instantiateStreaming(fetch("demo.wasm"), go.importObject).then((result) => {
				go.run(result.instance);
			});
		</script>
	</head>
	<body></body>
</html>`))
	case "/wasm_exec.js":
		http.ServeFile(w, r, wasmExecFile)
	case "/demo.wasm":
		http.ServeFile(w, r, *wasmFile)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	flag.Parse()

	if *tinygo {
		out, err := exec.Command("tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				log.Fatalf("%s: %s", exitError, exitError.Stderr)
			}
			log.Fatal(err)
		}
		out = bytes.TrimRight(out, "\n")
		wasmExecFile = filepath.Join(string(out), "targets", "wasm_exec.js")
	} else {
		wasmExecFile = filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
	}

	for _, path := range []string{wasmExecFile, *wasmFile} {
		if _, err := os.Stat(path); err != nil {
			log.Fatal(err)
		}
	}

	s := http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", *port),
		BaseContext: func(l net.Listener) context.Context {
			log.Printf("Serving on http://%s", l.Addr())
			return context.Background()
		},
		Handler: http.HandlerFunc(handler),
	}

	log.Fatal(s.ListenAndServe())
}
