package session

import (
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
)

type (
	redisStore struct {
		*redistore.RediStore
	}
)

func NewRedisStore(size int, network, address, password string, keyPairs ...[]byte) (Store, error) {
	store, err := redistore.NewRediStore(size, network, address, password, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &redisStore{store}, nil
}

func (c *redisStore) Options(options Options) {
	c.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
