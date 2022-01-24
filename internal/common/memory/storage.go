package memory

import (
	"sync"
	"time"
)

type Storage struct {
	mux        sync.RWMutex
	db         map[string]entry
	gcInterval time.Duration
	done       chan struct{}
}

type entry struct {
	expiry uint32
	data   []byte
}

func New() *Storage {
	store := &Storage{
		db:         make(map[string]entry),
		gcInterval: 10 * time.Second,
		done:       make(chan struct{}),
	}
	go store.gc()
	return store
}

func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	s.mux.RLock()
	v, ok := s.db[key]
	s.mux.RUnlock()
	if !ok || v.expiry != 0 && v.expiry <= uint32(time.Now().Unix()) {
		return nil, nil
	}

	return v.data, nil
}

func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	var expire uint32
	if exp != 0 {
		expire = uint32(time.Now().Add(exp).Unix())
	}

	s.mux.Lock()
	s.db[key] = entry{expire, val}
	s.mux.Unlock()
	return nil
}

func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	s.mux.Lock()
	delete(s.db, key)
	s.mux.Unlock()
	return nil
}

func (s *Storage) Reset() error {
	s.mux.Lock()
	s.db = make(map[string]entry)
	s.mux.Unlock()
	return nil
}

func (s *Storage) Close() error {
	s.done <- struct{}{}
	return nil
}

func (s *Storage) gc() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			now := uint32(t.Unix())
			s.mux.Lock()
			for id, v := range s.db {
				if v.expiry != 0 && v.expiry < now {
					delete(s.db, id)
				}
			}
			s.mux.Unlock()
		}
	}
}
