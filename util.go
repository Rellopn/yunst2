package yunst2

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"log"
	"os"
	"strings"
)

var (
	PFX_PATH    = ""
	PFX_PWD     = ""
	NO_PFX_PATH = errors.New("can not find PFX_PATH")
	NO_PFX_PWD  = errors.New("can not find PFX_PWD")
	privateKey  interface{}
	certificate *x509.Certificate
)

func SetPfxPath(pfxPath string) {
	PFX_PATH = pfxPath
}

func SetPfxPwd(pfxPwd string) {
	PFX_PWD = pfxPwd
}
func GetPair() {
	bytes, err := os.ReadFile(PFX_PATH)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, certificate, err = pkcs12.Decode(bytes, PFX_PWD)
	if err != nil {
		log.Fatal(err)
	}
}
func Sign(params string) (string, error) {
	if PFX_PATH == "" {
		return "", NO_PFX_PATH
	}
	if PFX_PWD == "" {
		return "", NO_PFX_PWD
	}
	h := md5.New()
	h.Write([]byte(params))
	encodeString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	s1 := sha1.New()
	s1.Write([]byte(encodeString))
	sum := s1.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, sum)
	if err != nil {
		log.Fatal(err)
	}
	singStr := base64.StdEncoding.EncodeToString(sig)
	return singStr, nil
}
func EncryptionSI(information string) (string, error) {
	if PFX_PATH == "" {
		return "", NO_PFX_PATH
	}
	if PFX_PWD == "" {
		return "", NO_PFX_PWD
	}
	block, _ := pem.Decode(caCert)
	var cert *x509.Certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	sig, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(information))
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(sig)), nil
}

func DecryptionSI(information string) (string, error) {
	if PFX_PATH == "" {
		return "", NO_PFX_PATH
	}
	if PFX_PWD == "" {
		return "", NO_PFX_PWD
	}
	deBytes, err := hex.DecodeString(information)
	if err != nil {
		return "", err
	}
	decodeBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), deBytes)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}

func VerifySign(signSource string, sign string) error {
	h := md5.New()
	h.Write([]byte(signSource))
	encodesignSource := base64.StdEncoding.EncodeToString(h.Sum(nil))
	encodesignSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(caCert)
	var cert *x509.Certificate
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	s1 := sha1.New()
	s1.Write([]byte(encodesignSource))
	sum := s1.Sum(nil)
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA1, sum, encodesignSign)
	return err
}
