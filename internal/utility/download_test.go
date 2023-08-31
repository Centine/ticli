package utility

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	ticlicrypto "github.com/centine/ticli/crypto"
)

func TestDownloadFileSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}))
	defer server.Close()

	path := "testfile.txt"
	err := downloadFile(path, server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer os.Remove(path)

	// Additional checks on the downloaded file content can be performed here.
}

func TestDownloadFileFail(t *testing.T) {
	err := downloadFile("testfile.txt", "http://invalid-url")
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

var privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

func TestVerifySignatureSuccess(t *testing.T) {
	message := []byte("Test Message")
	hashed := sha256.Sum256(message)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	// Deriving the public key from the private key
	publicKey := privateKey.Public()

	// Encode the public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("Error marshaling public key: %v", err)
	}
	// WARNING: HAS SIDE EFFECTS, OVERWRITES GLOBAL VARIABLE
	ticlicrypto.PublicKeyPem = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	err = verifySignature(message, signature)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestVerifySignatureFail(t *testing.T) {
	message := []byte("Test Message")
	wrongSignature := []byte("Wrong Signature")

	err := verifySignature(message, wrongSignature)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}
