package cache

import (
	"github.com/ashermp9/fiskaly-test-task/internal/domain"
	"github.com/ashermp9/fiskaly-test-task/pkg/cache"
)

type InMemoryStorage struct {
	DeviceCache *cache.Cache[string, domain.SignatureDevice]
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		DeviceCache: cache.NewCache[string, domain.SignatureDevice](),
	}
}
