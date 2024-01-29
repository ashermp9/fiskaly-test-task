package storage

import (
	"context"
	"fmt"

	"github.com/ashermp9/fiskaly-test-task/internal/domain"
)

// AddDevice adds a new signature device to the storage.
func (s *Storage) AddDevice(_ context.Context, device domain.SignatureDevice) {
	s.cache.DeviceCache.Set(device.ID, device)
}

// GetDevice retrieves a signature device from the storage.
func (s *Storage) GetDevice(_ context.Context, id string) (domain.SignatureDevice, error) {
	device, found := s.cache.DeviceCache.Get(id)
	if !found {
		return domain.SignatureDevice{}, fmt.Errorf("device not found")
	}
	return device, nil
}

func (s *Storage) LockDevice(_ context.Context, deviceID string) {
	s.cache.DeviceCache.Lock(deviceID)
}

func (s *Storage) UnlockDevice(_ context.Context, deviceID string) {
	s.cache.DeviceCache.Unlock(deviceID)
}
