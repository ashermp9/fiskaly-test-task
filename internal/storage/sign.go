package storage

import (
	"context"
	"fmt"

	"github.com/ashermp9/fiskaly-test-task/internal/domain"
)

// GenerateKeys uses CryptoManager to create key pairs.
func (s *Storage) GenerateKeys(ctx context.Context, algorithm domain.Algorithm) ([]byte, []byte, error) {
	generator, err := s.cryptoMgr.GetGenerator(algorithm)
	if err != nil {
		return nil, nil, err
	}
	return generator.GenerateBytes()
}

// SignTransaction uses CryptoManager to sign data.
func (s *Storage) SignTransaction(ctx context.Context, deviceID string, data []byte) ([]byte, error) {
	device, found := s.cache.DeviceCache.Get(deviceID)
	if !found {
		return nil, fmt.Errorf("device not found")
	}

	signer, err := s.cryptoMgr.GetSigner(device.Algorithm, device.PrivateKey)
	if err != nil {
		return nil, err
	}

	return signer.Sign(data)
}
