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
	mux.HandleFunc("/validate", pkg.Validate)
	mux.HandleFunc("/mutate", pkg.Mutate)

	log.Println("server listen on ':8000'...")

	if certFile != "" && keyFile != "" {
		log.Fatal(http.ListenAndServeTLS(":8000", certFile, keyFile, mux))
	} else {
		log.Fatal(http.ListenAndServe(":8000", mux))
	}
}
