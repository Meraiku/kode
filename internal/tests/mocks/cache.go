package mocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/meraiku/kode/internal/token"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	Cache map[string]any
	mock.Mock
}

func NewMockCache() *MockCache {
	return &MockCache{
		Cache: map[string]any{},
	}
}

func (c *MockCache) SetTokens(id string, tokens *token.Tokens, ctx context.Context) error {

	json, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	c.Cache[id] = json
	return nil
}

func (c *MockCache) GetTokens(id string, ctx context.Context) (*token.Tokens, error) {
	result, ok := c.Cache[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	switch v := result.(type) {
	case token.Tokens:
		return &v, nil
	default:
		log.Printf("%T", v)
		return nil, fmt.Errorf("%T", v)
	}

}

func (c *MockCache) DeleteTokens(id string, ctx context.Context) error {
	delete(c.Cache, id)
	return nil
}
