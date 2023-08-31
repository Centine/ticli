package utility

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/csv"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	ticlicrypto "github.com/centine/ticli/crypto"
	"github.com/centine/ticli/internal/config"
)

// var publicKeyPem = []byte(`
// -----BEGIN PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8TEyfM4IPlwzvlnxIpWa
// Vcw+Q8/7IIOC1oAMVApM/l+S2V9P/JjZE+KUB5ETK9nA5dzl4BEJRU37+i9t0MUO
// G+aJUDpcqmF1lBX6zOA6p44jb3kVt2vm7QLxnzFv8946fmsYc0sCei7HihC0t7qA
// PExtZ6S5xy1ojCnzHLaKK4wSe/tRmML2mWBYRap9JLPQ9mmqaz26ntYcXPsHoI4a
// DfP8Zuw5sH3tdy6Hk6GkKHIGxputMmRCdexgD1MaqJ9wtUNUYHM2K2Bvy7QQ998u
// xdqYnZfuKvfpN/QgVkmlP81SYyyO3DFkPio5MIx11UfIMoK8hKILTGDYegveMS8D
// KQIDAQAB
// -----END PUBLIC KEY-----
// `)

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
	log.Println("Downloaded file to", filepath)
	return err
}

func verifySignature(fileData []byte, signature []byte) error {
	block, _ := pem.Decode(ticlicrypto.PublicKeyPem)
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
	result := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	log.Printf("Signature verification result: %v\n", result) // nil is good
	return result
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

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DownloadAndVerify(ticliCtx *config.TicliContext) {
	scriptbundleLocalPath := ticliCtx.TicliDir + "/scriptbundles"
	bundle_url := "http://localhost:12845/scriptbundles_latest.zip" // FIXME
	sig_url := bundle_url + ".sig"
	script_path := scriptbundleLocalPath + "/scriptbundles_latest.zip"
	sig_path := scriptbundleLocalPath + "/scriptbundles_latest.zip.sig"

	err := downloadFile(script_path, bundle_url)
	if err != nil {
		panic(err)
	}
	err = downloadFile(sig_path, sig_url)
	if err != nil {
		panic(err)
	}

	scriptData, err := os.ReadFile(script_path)
	if err != nil {
		panic(err)
	}

	signature, err := os.ReadFile(sig_path)
	if err != nil {
		panic(err)
	}

	if err := verifySignature(scriptData, signature); err != nil {
		panic("Signature verification failed: " + err.Error())
	}

	err = Unzip(script_path, scriptbundleLocalPath)
	if err != nil {
		fmt.Println("Error unzipping file:", err)
	} else {
		log.Println("Unzipped successfully to", scriptbundleLocalPath)
	}

}
