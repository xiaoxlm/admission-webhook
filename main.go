package main

import (
	"github.com/xiaoxlm/admission-webhook/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	TLSCertName = "tls.crt"
	TLSKeyName  = "tls.key"
)

func main() {
	var (
		certFile string
		keyFile  string
	)
	if certDir := os.Getenv("CERT_DIR"); certDir != "" {
		certFile = filepath.Join(certDir, TLSCertName)
		keyFile = filepath.Join(certDir, TLSKeyName)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		return
	})
	mux.HandleFunc("/validate", pkg.Validate)
	mux.HandleFunc("/mutate", pkg.Mutate)

	log.Println("server listen on ':80'...")

	if certFile != "" && keyFile != "" {
		log.Fatal(http.ListenAndServeTLS(":80", certFile, keyFile, mux))
	} else {
		log.Fatal(http.ListenAndServe(":80", mux))
	}
}
