package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilovepitsa/happy/auth/api/sessions"
	"github.com/ilovepitsa/happy/backend/pkg/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	ErrNoAuth      = errors.New("no auth")
	ErrAuthService = errors.New("auth service error")
)

type (
	UserInterface interface {
		GetID() uint32
	}

	Session struct {
		SessionId string
		UserId    uint32
		Ttl       int64
	}

	SessionClient struct {
		client sessions.AuthCheckerClient
		cfg    *config.Config
	}
	ctxKey int
)

const SessionKey ctxKey = 1

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}

	return sess, nil
}

func NewAuthClient(conn grpc.ClientConnInterface, cfg *config.Config) *SessionClient {
	return &SessionClient{
		client: sessions.NewAuthCheckerClient(conn),
		cfg:    cfg,
	}
}

func (s *SessionClient) Check(c *gin.Context) (*Session, error) {
	sessionCookie, err := c.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Debug("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{}
	res, err := s.client.Check(c, &sessions.SessionID{ID: sessionCookie})

	if err != nil {
		log.Debug("Error auth service")
		return nil, ErrAuthService
	}

	sess.SessionId = sessionCookie
	sess.UserId = res.UserID

	return sess, nil
}

func (s *SessionClient) Create(c *gin.Context, user UserInterface) error {
	sessId, err := s.client.Create(context.TODO(), &sessions.Session{UserID: user.GetID()})
	if err != nil {
		log.Debug("error create session")
		return err
	}

	c.SetCookie("session_id", sessId.ID, int(sessId.Ttl), "/", s.cfg.Net.Host, false, false)

	return nil
}

func (s *SessionClient) Delete(c *gin.Context) error {
	_, err := SessionFromContext(c)
	if err != nil {
		log.Debug("error delete session")
		return err
	}

	c.SetCookie("session_id", "", -1, "/", s.cfg.Net.Host, false, false)

	return nil
}
