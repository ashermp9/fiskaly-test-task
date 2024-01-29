package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type KeyGenerator interface {
	GenerateBytes() (publicKey []byte, privateKey []byte, err error)
}

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// GenerateBytes generates a new ECCKeyPair and returns encoded keys.
func (g *ECCGenerator) GenerateBytes() ([]byte, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE_KEY", Bytes: privateKeyBytes})
	encodedPublic := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC_KEY", Bytes: publicKeyBytes})

	return encodedPublic, encodedPrivate, nil
}

// GenerateBytes generates a new RSAKeyPair and returns encoded keys.
func (g *RSAGenerator) GenerateBytes() ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	encodedPrivate := pem.EncodeToMemory(&pem.Block{Type: "RSA_PRIVATE_KEY", Bytes: privateKeyBytes})
	encodedPublic := pem.EncodeToMemory(&pem.Block{Type: "RSA_PUBLIC_KEY", Bytes: publicKeyBytes})

	return encodedPublic, encodedPrivate, nil
}
