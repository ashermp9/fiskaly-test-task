package crypto

import (
	"fmt"
	"sync"

	"github.com/ashermp9/fiskaly-test-task/internal/domain"
	"github.com/ashermp9/fiskaly-test-task/pkg/crypto"
)

// CryptoManager manages cryptographic generators and signers.
type CryptoManager struct {
	generators map[domain.Algorithm]crypto.KeyGenerator
	signers    map[domain.Algorithm]crypto.Signer
	mu         sync.RWMutex
}

// NewCryptoManager creates a new instance of CryptoManager.
func NewCryptoManager() *CryptoManager {
	return &CryptoManager{
		generators: make(map[domain.Algorithm]crypto.KeyGenerator),
		signers:    make(map[domain.Algorithm]crypto.Signer),
	}
}

// GetGenerator retrieves a key generator based on the specified algorithm.
func (m *CryptoManager) GetGenerator(algorithm domain.Algorithm) (crypto.KeyGenerator, error) {
	m.mu.RLock()
	generator, exists := m.generators[algorithm]
	m.mu.RUnlock()

	if !exists {
		switch algorithm {
		case domain.AlgorithmRSA:
			generator = &crypto.RSAGenerator{}
		case domain.AlgorithmECC:
			generator = &crypto.ECCGenerator{}
		default:
			return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
		}

		m.mu.Lock()
		m.generators[algorithm] = generator
		m.mu.Unlock()
	}

	return generator, nil
}

// GetSigner retrieves a signer based on the specified algorithm and private key.
func (m *CryptoManager) GetSigner(algorithm domain.Algorithm, privateKey []byte) (crypto.Signer, error) {
	m.mu.RLock()
	signer, exists := m.signers[algorithm]
	m.mu.RUnlock()

	if !exists {
		switch algorithm {
		case domain.AlgorithmRSA:
			rsaKeyPair, err := crypto.NewRSAMarshaler().Unmarshal(privateKey)
			if err != nil {
				return nil, err
			}
			signer = crypto.NewRSASigner(rsaKeyPair.Private)
		case domain.AlgorithmECC:
			eccKeyPair, err := crypto.NewECCMarshaler().Decode(privateKey)
			if err != nil {
				return nil, err
			}
			signer = crypto.NewECDSASigner(eccKeyPair.Private)
		default:
			return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
		}

		m.mu.Lock()
		m.signers[algorithm] = signer
		m.mu.Unlock()
	}

	return signer, nil
}
