package memory

import (
	"errors"
	"github.com/SemmiDev/blog/internal/common/memory"
	"github.com/SemmiDev/blog/internal/user/query"
)

type TokenQueryMemory struct {
	DB *memory.Storage
}

func NewTokenQueryMemory(DB *memory.Storage) *TokenQueryMemory {
	return &TokenQueryMemory{DB: DB}
}

func (s *TokenQueryMemory) Find(key string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		data, _ := s.DB.Get(key)
		if data == nil {
			result <- query.Result{Error: errors.New("data not found")}
		}

		result <- query.Result{Result: data}
		close(result)
	}()

	return result
}
