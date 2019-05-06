package main

import (
	"net/http"
	"strings"
)

const (
	wasmSuffix      = ".wasm"
	wasmContentType = "application/wasm"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, wasmSuffix) {
			res.Header().Add("Content-Type", wasmContentType)
		}

		http.FileServer(http.Dir(".")).ServeHTTP(res, req)
	})
	http.ListenAndServe(":8000", nil)
}
