package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/meraiku/kode/internal/token"
)

func (c *Redis) SetTokens(id string, tokens *token.Tokens, ctx context.Context) error {

	json, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	_, err = c.cache.Set(ctx, id, json, time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) GetTokens(id string, ctx context.Context) (*token.Tokens, error) {
	result, err := c.cache.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	tokens := &token.Tokens{}

	if err := json.Unmarshal([]byte(result), tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (c *Redis) DeleteTokens(id string, ctx context.Context) error {
	_, err := c.cache.Del(ctx, id).Result()
	if err != nil {
		return err
	}
	return nil
}
