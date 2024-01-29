package storage

import (
	"github.com/ashermp9/fiskaly-test-task/internal/adapters/cache"
	"github.com/ashermp9/fiskaly-test-task/internal/adapters/crypto"
)

type Storage struct {
	cache     *cache.InMemoryStorage
	cryptoMgr *crypto.CryptoManager
}

func NewStorage() *Storage {
	return &Storage{
		cache:     cache.NewInMemoryStorage(),
		cryptoMgr: crypto.NewCryptoManager(),
	}
}
