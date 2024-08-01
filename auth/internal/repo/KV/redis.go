package kv

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ilovepitsa/happy/auth/api/sessions"
	"github.com/ilovepitsa/happy/auth/pkg/config"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

const (
	defaultConnAttempts = 5
	defaultConnTimeout  = time.Second
)

type Redis struct {
	mtx          sync.RWMutex
	client       *redis.Client
	connAttempts int
	ttl          time.Duration
	connTimeout  time.Duration
}

func New(config *config.Config) (*Redis, error) {
	r := &Redis{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.R.Host, config.R.Port),
		Password: config.R.Password,
		DB:       0,
	})
	var err error
	for r.connAttempts > 0 {

		if err = r.client.Ping(context.Background()).Err(); err == nil {
			break
		}

		log.Info("Attempting to connect redis. Attempts left: ", r.connAttempts)
		time.Sleep(r.connTimeout)
		r.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis new %w", err)
	}
	return r, nil
}

func (r *Redis) GetTTL() time.Duration {
	return r.ttl
}

func (r *Redis) Store(key *sessions.SessionID, value *sessions.Session) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	err := r.client.Set(context.TODO(), key.ID, value.UserID, r.ttl).Err()
	return err
}
func (r *Redis) GetSess(key *sessions.SessionID) (*sessions.Session, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	res, err := r.client.Get(context.TODO(), key.ID).Result()
	if err != nil {
		return nil, err
	}

	usrId, err := strconv.ParseUint(res, 10, 32)
	if err != nil {
		return nil, err
	}

	return &sessions.Session{UserID: uint32(usrId)}, nil
}

func (r *Redis) Del(key *sessions.SessionID) error {
	_, err := r.client.Del(context.TODO(), key.ID).Result()
	return err
}

func (r *Redis) Close() {
	if r.client != nil {
		r.client.Close()
	}
}
