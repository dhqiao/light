package chain

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"encoding/base64"
)

// test file read pem file decode
func GetCertFileInfo(path string) *x509.Certificate {
	creatorByte, err := ioutil.ReadFile(path)
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	uname:=cert.Subject.CommonName
	fmt.Println("Name:"+uname)
	decodeExt(*cert)
	return cert
}

func decodeExt(cert x509.Certificate)  {
	for i := 0; i < len(cert.Extensions); i++ {
		value := cert.Extensions[i].Value
		stringR, err := base64.StdEncoding.DecodeString(string(value))
		if err != nil {}
		fmt.Println("-------------", string(value))
		fmt.Println("..............", string(stringR[:]))
		fmt.Println(string(stringR[:]))
	}
}