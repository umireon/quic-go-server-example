package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go/http3"
)

func isCertExists() bool {
	fileCert, errCert := os.Open("cert.pem")
	fileCert.Close()
	if errors.Is(errCert, os.ErrNotExist) {
		return false
	}

	fileKey, errKey := os.Open("key.pem")
	fileKey.Close()
	if errors.Is(errKey, os.ErrNotExist) {
		return false
	}

	return true
}

func generateTLSCert() {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	keyOut, _ := os.Create("key.pem")
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certOut, _ := os.Create("cert.pem")
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
}

func hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	if !isCertExists() {
		generateTLSCert()
	}
	http3.ListenAndServe("0.0.0.0:4242", "cert.pem", "key.pem", http.HandlerFunc(hello))
}
