package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// ReadPublicKeyFromFile reads a PEM-encoded RSA public key from the file located at the given path
func ReadPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	bz, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	pubPem, _ := pem.Decode(bz)
	if pubPem == nil {
		return nil, fmt.Errorf(
			"rsa public key not in pem format" +
				"Use `ssh-keygen -f id_rsa.pub -e -m pem > id_rsa.pem` to generate the pem encoding of your RSA public key",
		)
	}

	return x509.ParsePKCS1PublicKey(pubPem.Bytes)
}
