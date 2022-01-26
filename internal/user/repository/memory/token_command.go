package memory

import (
	"github.com/SemmiDev/blog/internal/common/memory"
	"time"
)

type TokenCommandMemory struct {
	DB *memory.Storage
}

func NewTokenCommandMemory(DB *memory.Storage) *TokenCommandMemory {
	return &TokenCommandMemory{DB: DB}
}

func (t *TokenCommandMemory) Set(key string, val []byte, exp time.Duration) <-chan error {
	result := make(chan error)

	go func() {
		err := t.DB.Set(key, val, exp)
		if err != nil {
			result <- nil
		}

		result <- nil
		close(result)
	}()

	return result
}

func (t *TokenCommandMemory) Delete(key string) <-chan error {
	result := make(chan error)

	go func() {
		err := t.DB.Delete(key)
		if err != nil {
			result <- nil
		}

		result <- nil
		close(result)
	}()

	return result
}
