package utility

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/csv"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var publicKeyPem = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8TEyfM4IPlwzvlnxIpWa
Vcw+Q8/7IIOC1oAMVApM/l+S2V9P/JjZE+KUB5ETK9nA5dzl4BEJRU37+i9t0MUO
G+aJUDpcqmF1lBX6zOA6p44jb3kVt2vm7QLxnzFv8946fmsYc0sCei7HihC0t7qA
PExtZ6S5xy1ojCnzHLaKK4wSe/tRmML2mWBYRap9JLPQ9mmqaz26ntYcXPsHoI4a
DfP8Zuw5sH3tdy6Hk6GkKHIGxputMmRCdexgD1MaqJ9wtUNUYHM2K2Bvy7QQ998u
xdqYnZfuKvfpN/QgVkmlP81SYyyO3DFkPio5MIx11UfIMoK8hKILTGDYegveMS8D
KQIDAQAB
-----END PUBLIC KEY-----
`)

func downloadFile(filepath string, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func verifySignature(fileData []byte, signature []byte) error {
	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return fmt.Errorf("failed to decode public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key has wrong type")
	}

	hashed := sha256.Sum256(fileData)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
}

func executeScript(scriptPath string) ([][]string, error) {
	cmd := exec.Command("/bin/bash", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func main() {
	url := "https://x.y.z/cli/linuxchecks_latest.sh"
	scriptPath := "linuxchecks_latest.sh"
	signaturePath := "linuxchecks_latest.sig" // Assuming the signature is available at this path

	err := downloadFile(scriptPath, url)
	if err != nil {
		panic(err)
	}

	scriptData, err := os.ReadFile(scriptPath)
	if err != nil {
		panic(err)
	}

	signature, err := os.ReadFile(signaturePath)
	if err != nil {
		panic(err)
	}

	if err := verifySignature(scriptData, signature); err != nil {
		panic("Signature verification failed: " + err.Error())
	}

	records, err := executeScript(scriptPath)
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		fmt.Println(record)
	}
}
