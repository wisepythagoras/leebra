package js

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

func ECDSAEncodePrivate(pk *ecdsa.PrivateKey) string {
	x509EncodedPriv, _ := x509.MarshalECPrivateKey(pk)
	pemPrivBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509EncodedPriv,
	}

	return string(pem.EncodeToMemory(pemPrivBlock))
}

func ECDSAEncodePublic(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemPubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509EncodedPub,
	}

	return string(pem.EncodeToMemory(pemPubBlock))
}

func ECDSADecodePrivate(pemPriv string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemPriv))
	x509Encoded := block.Bytes
	return x509.ParseECPrivateKey(x509Encoded)
}

func ECDSADecodePublic(pemPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)

	if err != nil {
		return nil, err
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}
