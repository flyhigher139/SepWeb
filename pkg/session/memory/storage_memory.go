package memory

import (
	"context"
	"errors"
	"github.com/igevin/sepweb/pkg/session"
	"github.com/patrickmn/go-cache"
	"time"
)

type Store struct {
	c          *cache.Cache
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		c:          cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	s.c.Set(sess.ID(), sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	sess, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	s.c.Set(sess.ID(), sess, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.c.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	sess, ok := s.c.Get(id)
	if !ok {
		return nil, errors.New("session is no found")
	}
	return sess.(*memorySession), nil
}

type memorySession struct {
	id         string
	data       map[string]string
	expiration time.Duration
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.data[key]; !ok {
		return "", errors.New("找不到这个key")
	} else {
		return val, nil
	}

}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.data[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.id
}
