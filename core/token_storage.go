package core

import (
	"encoding/json"
	"log"
)

type TokenStorage interface {
	Save(repo string, token *Token)
	Read(repo string) (*Token, error)
}

type InMemoryTokenStorage struct {
	store map[string][]byte
}

func StorageSetup() TokenStorage {
	return &InMemoryTokenStorage{
		store: make(map[string][]byte),
	}
}

func (s *InMemoryTokenStorage) Save(repo string, token *Token) {
	j, _ := json.Marshal(token)
	s.store[repo] = j
}

func (s *InMemoryTokenStorage) Read(repo string) (*Token, error) {
	var t Token
	err := json.Unmarshal(s.store[repo], &t)
	if err != nil {
		log.Printf("Faile to read session cookie: %v", err)
		return nil, err
	}
	return &t, nil
}
