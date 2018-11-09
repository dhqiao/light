package chain

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"encoding/base64"
)

// test file

func GetCertFileInfo(path string) *x509.Certificate {
	//keyPath := "service/crypto-config/peerOrganizations/member1.example.com/users/Admin@member1.example.com/tls/client.key"

	var creatorByte []byte
	var err error
	creatorByte, err = ioutil.ReadFile(path)


	//creatorByte,_:= stub.GetCreator()
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
		if err != nil {

		}
		fmt.Println(string(stringR[:]))
	}
}