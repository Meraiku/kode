package cache

import (
	"context"

	"github.com/meraiku/kode/internal/token"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	SetTokens(id string, tokens *token.Tokens, ctx context.Context) error
	GetTokens(id string, ctx context.Context) (*token.Tokens, error)
	DeleteTokens(refresh string, ctx context.Context) error
}

type Redis struct {
	cache *redis.Client
}

func NewCache(cache *redis.Client) Cache {
	return &Redis{
		cache: cache,
	}
}
