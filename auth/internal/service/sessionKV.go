package service

import (
	"math/rand"
	"time"

	"github.com/ilovepitsa/happy/auth/api/sessions"
)

const sessKeyLen = 10

type CacheRepo interface {
	Store(key *sessions.SessionID, value *sessions.Session) error
	GetSess(key *sessions.SessionID) (*sessions.Session, error)
	Del(key *sessions.SessionID) error
	GetTTL() time.Duration
}

type KVSessionManager struct {
	repo CacheRepo
}

func (sm *KVSessionManager) Check(sessId *sessions.SessionID) (*sessions.Session, error) {
	return sm.repo.GetSess(sessId)
}
func (sm *KVSessionManager) Create(sess *sessions.Session) (*sessions.SessionID, error) {
	id := &sessions.SessionID{
		ID:  randStringRunes(sessKeyLen),
		Ttl: sm.repo.GetTTL().Nanoseconds()}
	err := sm.repo.Store(id, sess)
	if err != nil {
		return nil, err
	}

	return id, nil
}
func (sm *KVSessionManager) Delete(sessId *sessions.SessionID) error {
	return sm.repo.Del(sessId)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewKVSessManager(repo CacheRepo) *KVSessionManager {
	return &KVSessionManager{
		repo: repo,
	}
}
