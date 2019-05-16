package yunst2

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"log"
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
func getPair() {
	bytes, err := ioutil.ReadFile(PFX_PATH)
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
