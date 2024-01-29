package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{PrivateKey: privateKey}
}

func (s *RSASigner) Sign(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

type ECDSASigner struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewECDSASigner(privateKey *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{PrivateKey: privateKey}
}

func (signer *ECDSASigner) Sign(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, signer.PrivateKey, hashed[:])
	if err != nil {
		return nil, err
	}

	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}
